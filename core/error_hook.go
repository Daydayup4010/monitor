package core

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"
	"uu/config"
	"uu/utils"

	"github.com/sirupsen/logrus"
)

// ErrorEmailHook 错误邮件通知钩子
type ErrorEmailHook struct {
	emailService *utils.EmailService
	recipients   []string
	minLevel     logrus.Level
	rateLimit    int           // 每分钟最大邮件数
	cooldown     time.Duration // 相同错误冷却时间
	batchWindow  time.Duration // 批量发送窗口

	mu            sync.Mutex
	sentCount     int                  // 当前分钟已发送数量
	lastReset     time.Time            // 上次重置时间
	errorCache    map[string]time.Time // 错误去重缓存 (hash -> 上次发送时间)
	pendingErrors []errorEntry         // 待发送的错误（用于批量发送）
	batchTicker   *time.Ticker
	stopChan      chan struct{}
}

type errorEntry struct {
	Level     logrus.Level
	Message   string
	File      string
	Function  string
	Line      int
	Timestamp time.Time
	Fields    map[string]interface{}
}

// NewErrorEmailHook 创建错误邮件通知钩子
func NewErrorEmailHook(emailService *utils.EmailService, cfg *config.ErrorAlert) *ErrorEmailHook {
	if emailService == nil || cfg == nil || !cfg.Enabled {
		return nil
	}

	minLevel := logrus.ErrorLevel
	switch strings.ToLower(cfg.MinLevel) {
	case "error":
		minLevel = logrus.ErrorLevel
	case "fatal":
		minLevel = logrus.FatalLevel
	case "panic":
		minLevel = logrus.PanicLevel
	}

	rateLimit := cfg.RateLimit
	if rateLimit <= 0 {
		rateLimit = 10 // 默认每分钟最多10封
	}

	cooldown := time.Duration(cfg.Cooldown) * time.Second
	if cooldown <= 0 {
		cooldown = 5 * time.Minute // 默认5分钟冷却
	}

	batchWindow := time.Duration(cfg.BatchWindow) * time.Second

	hook := &ErrorEmailHook{
		emailService:  emailService,
		recipients:    cfg.Recipients,
		minLevel:      minLevel,
		rateLimit:     rateLimit,
		cooldown:      cooldown,
		batchWindow:   batchWindow,
		lastReset:     time.Now(),
		errorCache:    make(map[string]time.Time),
		pendingErrors: make([]errorEntry, 0),
		stopChan:      make(chan struct{}),
	}

	// 如果配置了批量发送窗口，启动定时发送
	if batchWindow > 0 {
		hook.batchTicker = time.NewTicker(batchWindow)
		go hook.batchSender()
	}

	// 定期清理过期的缓存
	go hook.cleanupCache()

	return hook
}

// Levels 返回钩子触发的日志级别
func (h *ErrorEmailHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)
	for _, level := range logrus.AllLevels {
		if level <= h.minLevel {
			levels = append(levels, level)
		}
	}
	return levels
}

// Fire 当日志触发时调用
func (h *ErrorEmailHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查速率限制
	now := time.Now()
	if now.Sub(h.lastReset) >= time.Minute {
		h.sentCount = 0
		h.lastReset = now
	}

	if h.sentCount >= h.rateLimit {
		return nil // 达到速率限制，跳过
	}

	// 生成错误指纹用于去重
	fingerprint := h.generateFingerprint(entry)

	// 检查是否在冷却期内
	if lastSent, exists := h.errorCache[fingerprint]; exists {
		if now.Sub(lastSent) < h.cooldown {
			return nil // 在冷却期内，跳过
		}
	}

	// 构建错误条目
	errEntry := errorEntry{
		Level:     entry.Level,
		Message:   entry.Message,
		Timestamp: entry.Time,
		Fields:    make(map[string]interface{}),
	}

	if entry.HasCaller() {
		errEntry.File = entry.Caller.File
		errEntry.Function = entry.Caller.Function
		errEntry.Line = entry.Caller.Line
	}

	for k, v := range entry.Data {
		errEntry.Fields[k] = v
	}

	// 更新缓存
	h.errorCache[fingerprint] = now

	if h.batchWindow > 0 {
		// 批量模式：添加到待发送队列
		h.pendingErrors = append(h.pendingErrors, errEntry)
	} else {
		// 立即发送模式
		go h.sendEmail([]errorEntry{errEntry})
		h.sentCount++
	}

	return nil
}

// generateFingerprint 生成错误指纹用于去重
func (h *ErrorEmailHook) generateFingerprint(entry *logrus.Entry) string {
	var parts []string
	parts = append(parts, entry.Level.String())
	parts = append(parts, entry.Message)

	if entry.HasCaller() {
		parts = append(parts, entry.Caller.File)
		parts = append(parts, fmt.Sprintf("%d", entry.Caller.Line))
	}

	data := strings.Join(parts, "|")
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// batchSender 批量发送协程
func (h *ErrorEmailHook) batchSender() {
	for {
		select {
		case <-h.batchTicker.C:
			h.mu.Lock()
			if len(h.pendingErrors) > 0 {
				errors := make([]errorEntry, len(h.pendingErrors))
				copy(errors, h.pendingErrors)
				h.pendingErrors = h.pendingErrors[:0]
				h.sentCount++
				h.mu.Unlock()
				h.sendEmail(errors)
			} else {
				h.mu.Unlock()
			}
		case <-h.stopChan:
			return
		}
	}
}

// cleanupCache 定期清理过期缓存
func (h *ErrorEmailHook) cleanupCache() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			now := time.Now()
			for key, t := range h.errorCache {
				if now.Sub(t) > h.cooldown*2 {
					delete(h.errorCache, key)
				}
			}
			h.mu.Unlock()
		case <-h.stopChan:
			return
		}
	}
}

// sendEmail 发送错误告警邮件
func (h *ErrorEmailHook) sendEmail(errors []errorEntry) {
	if len(errors) == 0 || len(h.recipients) == 0 {
		return
	}

	subject := fmt.Sprintf("【CS Goods】系统错误告警 - %s", time.Now().Format("2006-01-02 15:04:05"))

	var body strings.Builder
	body.WriteString(`
		<div style="font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #e74c3c; border-bottom: 2px solid #e74c3c; padding-bottom: 10px;">
				⚠️ CS Goods 系统错误告警
			</h2>
			<p style="color: #666;">检测到以下错误，请及时处理：</p>
	`)

	for i, err := range errors {
		levelColor := "#e74c3c"
		levelIcon := "❌"
		switch err.Level {
		case logrus.WarnLevel:
			levelColor = "#f39c12"
			levelIcon = "⚠️"
		case logrus.FatalLevel:
			levelColor = "#8e44ad"
			levelIcon = "💀"
		case logrus.PanicLevel:
			levelColor = "#2c3e50"
			levelIcon = "🔥"
		}

		body.WriteString(fmt.Sprintf(`
			<div style="background: #f8f9fa; border-left: 4px solid %s; padding: 15px; margin: 15px 0; border-radius: 4px;">
				<div style="margin-bottom: 10px;">
					<span style="background: %s; color: white; padding: 3px 8px; border-radius: 3px; font-size: 12px;">
						%s %s
					</span>
					<span style="color: #999; font-size: 12px; margin-left: 10px;">
						%s
					</span>
				</div>
				<div style="font-size: 14px; color: #333; margin: 10px 0;">
					<strong>错误信息：</strong> %s
				</div>
		`, levelColor, levelColor, levelIcon, strings.ToUpper(err.Level.String()), err.Timestamp.Format("2006-01-02 15:04:05"), err.Message))

		if err.File != "" {
			body.WriteString(fmt.Sprintf(`
				<div style="font-size: 12px; color: #666; margin: 5px 0;">
					<strong>位置：</strong> %s:%d
				</div>
			`, err.File, err.Line))
		}

		if err.Function != "" {
			body.WriteString(fmt.Sprintf(`
				<div style="font-size: 12px; color: #666; margin: 5px 0;">
					<strong>函数：</strong> %s
				</div>
			`, err.Function))
		}

		if len(err.Fields) > 0 {
			body.WriteString(`<div style="font-size: 12px; color: #666; margin: 5px 0;"><strong>附加信息：</strong><ul style="margin: 5px 0;">`)
			for k, v := range err.Fields {
				body.WriteString(fmt.Sprintf(`<li><code>%s</code>: %v</li>`, k, v))
			}
			body.WriteString(`</ul></div>`)
		}

		body.WriteString(`</div>`)

		// 限制单封邮件最多显示10条错误
		if i >= 9 && len(errors) > 10 {
			body.WriteString(fmt.Sprintf(`
				<p style="color: #999; text-align: center;">... 还有 %d 条错误未显示 ...</p>
			`, len(errors)-10))
			break
		}
	}

	body.WriteString(fmt.Sprintf(`
			<div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #999; font-size: 12px;">
				<p>服务器：%s</p>
				<p>发送时间：%s</p>
				<p>此邮件由系统自动发送，请勿直接回复。</p>
			</div>
		</div>
	`, config.CONFIG.Server.GetAddr(), time.Now().Format("2006-01-02 15:04:05")))

	// 发送邮件
	h.emailService.SendErrorAlert(h.recipients, subject, body.String())
}

// Stop 停止钩子
func (h *ErrorEmailHook) Stop() {
	close(h.stopChan)
	if h.batchTicker != nil {
		h.batchTicker.Stop()
	}
}
