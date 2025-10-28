package services

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"
	"uu/config"
	"uu/models"

	"golang.org/x/time/rate"
)

var UUMaxPageSize = 100
var BuffMaxPageSize = "80"
var RequestDelay = time.Second * 30 // 429后等30秒
var wg sync.WaitGroup
var taskBuff sync.Mutex
var taskUU sync.Mutex
var uuLimiter = rate.NewLimiter(3, 1)                           // 速率: 3 tokens/s, 突发容量: 1
var buffLimiter = rate.NewLimiter(rate.Every(3*time.Second), 1) // 每3秒1个请求

// 滑动窗口限流配置 - 充分利用Buff的10次/30秒限额
var (
	requestWindow       = 30 * time.Second // 窗口大小：30秒
	maxRequestsInWindow = 9                // 窗口内最多9个
	recentRequests      = make([]time.Time, 0, 10)
	requestMutex        sync.Mutex
)

func UpdateAllUUItems() {
	_, total, _ := GetUUItems(20, 1)
	totalPages := total/UUMaxPageSize + 1
	for page := 1; page <= totalPages+1; page++ {
		if err := uuLimiter.Wait(context.Background()); err != nil {
			config.Log.Errorf("wait limiter error: %v", err)
			return
		}

		items, _, err := GetUUItems(UUMaxPageSize, page)
		if err != nil {
			if isRateLimitError(err) {
				// 动态退避：遇到429时增加延迟
				handleRateLimitError()
				page-- // 重试当前页
				continue
			}
			config.Log.Errorf("request uu %d page fail: %v", page, err)
			continue
		}
		models.BatchAddUUItem(items)
		//config.Log.Infof("Full Update uu item pageName: %d, success", page)
	}
}

// 判断是否为速率限制错误
func isRateLimitError(err error) bool {
	return strings.Contains(err.Error(), "429") ||
		strings.Contains(err.Error(), "too many requests")
}

// 处理速率限制错误（动态退避）
func handleRateLimitError() {
	config.Log.Warnf("limits sleep %v", RequestDelay)
	time.Sleep(RequestDelay)
}

// 检查滑动窗口并控制请求频率
func checkRequestWindow() {
	requestMutex.Lock()
	defer requestMutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-requestWindow)

	// 清理窗口外的旧请求记录
	validRequests := make([]time.Time, 0)
	for _, t := range recentRequests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}
	recentRequests = validRequests

	// 如果窗口内已经有9个请求，需要等待最早的请求离开窗口
	if len(recentRequests) >= maxRequestsInWindow {
		oldestRequest := recentRequests[0]
		waitTime := requestWindow - now.Sub(oldestRequest) + time.Millisecond*100

		if waitTime > 0 {
			config.Log.Infof("窗口内已有%d个请求，等待%v后再发起", len(recentRequests), waitTime)
			requestMutex.Unlock() // 释放锁再sleep
			time.Sleep(waitTime)
			requestMutex.Lock() // 重新获取锁

			// 重新清理（sleep期间时间已过）
			now = time.Now()
			windowStart = now.Add(-requestWindow)
			validRequests = make([]time.Time, 0)
			for _, t := range recentRequests {
				if t.After(windowStart) {
					validRequests = append(validRequests, t)
				}
			}
			recentRequests = validRequests
		}
	}

	// 记录本次请求时间
	recentRequests = append(recentRequests, now)
	config.Log.Debugf("窗口内请求数: %d/%d", len(recentRequests), maxRequestsInWindow)
}

// 清空窗口记录（遇到429时调用）
func clearRequestWindow() {
	requestMutex.Lock()
	defer requestMutex.Unlock()
	recentRequests = make([]time.Time, 0, 10)
	config.Log.Info("清空请求窗口记录")
}

func UpdateAllBuffItems() {
	_, total, _ := GetBuffItems("20", "1")
	size, _ := strconv.Atoi(BuffMaxPageSize)
	pageNum := total/size + 1

	for i := 1; i <= pageNum+1; i++ {
		// 1. 先通过rate limiter（保证基础间隔3秒）
		if err := buffLimiter.Wait(context.Background()); err != nil {
			config.Log.Errorf("wait buff limiter error: %v", err)
			return
		}

		// 2. 再检查滑动窗口（确保30秒内不超过9次）
		checkRequestWindow()

		// 3. 发起请求
		items, _, err := GetBuffItems(BuffMaxPageSize, strconv.Itoa(i))

		if err != nil {
			if isRateLimitError(err) {
				config.Log.Warnf("request buff %d page 触发429", i)

				// 清空窗口记录，重新开始计数
				clearRequestWindow()

				// 动态退避：等待30秒
				handleRateLimitError()

				// 重试当前页
				i--
				continue
			}
			config.Log.Errorf("request buff %d page fail: %v", i, err)
			continue
		}

		config.Log.Infof("request buff %d page success", i)
		models.BatchAddBuffItem(items)
	}
}

//func UpdateFullData() int {
//	if !task.TryLock() {
//		config.Log.Info("full update running")
//		return utils.ErrCodeFullUpdateRunning
//	}
//	defer task.Unlock()
//	config.Log.Info("Start full update")
//	wg.Add(2)
//	go UpdateAllUUItems()
//	go UpdateAllBuffItems()
//	wg.Wait()
//	config.Log.Info("Full update completed")
//	return utils.SUCCESS
//}

func UpdateInventory() {
	uu := GetUUInventory()
	models.BatchAddUUInventory(uu)
	buff := GetBuffInventory()
	models.BatchAddBuffInventory(buff)
}

func UpdateUUFullData() {
	if !taskUU.TryLock() {
		config.Log.Info("uu full update running")
	}
	defer taskUU.Unlock()
	config.Log.Info("Start uu full update")
	UpdateAllUUItems()
	config.Log.Info("uu full update completed")
}

func UpdateBuffFullData() {
	if !taskBuff.TryLock() {
		config.Log.Info("buff full update running")
	}
	defer taskBuff.Unlock()
	config.Log.Info("Start buff full update")
	UpdateAllBuffItems()
	config.Log.Info("buff full update completed")
}
