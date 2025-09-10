package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
	"uu/config"

	rota "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 字体颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
}

// Format 实现Formatter接口里的Format(*Entry) ([]byte, error)
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	// 日志前缀
	log := config.CONFIG.Logger.Prefix
	// 自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		// 自定义log存放路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// 自定义输出格式
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", log, timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s\n", log, timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() {
	mLog := logrus.New()

	// 创建log目录
	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		return
	}

	// 设置日志文件路径和名称
	fileName := "app.log"
	filePath := filepath.Join(logDir, fileName)
	linkName := filepath.Join(logDir, "latest.log")

	// 创建日志文件轮转器
	logWriter, err := rota.New(
		filePath+"%Y%m%d.log",               // 日志文件名格式
		rota.WithMaxAge(7*24*time.Hour),     // 保留7天的日志
		rota.WithRotationTime(24*time.Hour), // 每天轮转一次
		rota.WithLinkName(linkName),         // 创建最新日志的软链接
	)
	if err != nil {
		fmt.Printf("创建日志轮转器失败: %v\n", err)
		mLog.SetOutput(os.Stdout) // 如果文件输出失败，回退到标准输出
	} else {
		// 根据配置决定输出方式
		var writers []io.Writer
		writers = append(writers, logWriter) // 总是输出到文件

		if config.CONFIG.Logger.LogInConsole {
			writers = append(writers, os.Stdout) // 如果配置允许，也输出到控制台
		}

		// 创建多重输出
		multiWriter := io.MultiWriter(writers...)
		mLog.SetOutput(multiWriter)

		// 为不同级别的日志设置文件输出钩子（使用自定义格式器）
		wireMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.DebugLevel: logWriter,
			logrus.FatalLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.PanicLevel: logWriter,
			logrus.ErrorLevel: logWriter,
		}

		// 为文件输出使用简化的格式（不带颜色）
		fileHook := lfshook.NewHook(wireMap, &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   true, // 文件输出不需要颜色
		})
		mLog.AddHook(fileHook)
	}

	mLog.SetReportCaller(config.CONFIG.Logger.ShowLine) // 开启返回函数名和行号
	mLog.SetFormatter(&LogFormatter{})                  // 设置自己定义的formatter（用于控制台输出）
	parseLevel, err := logrus.ParseLevel(config.CONFIG.Logger.Level)
	if err != nil {
		parseLevel = logrus.InfoLevel
	}
	mLog.SetLevel(parseLevel) // 设置level
	config.Log = mLog
}
