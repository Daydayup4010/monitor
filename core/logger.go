package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"uu/config"
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
	timestamp := entry.Time.Format("2006-01-02 15:01:05")
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
	mLog.SetOutput(os.Stdout)                           // 设置输出类型
	mLog.SetReportCaller(config.CONFIG.Logger.ShowLine) // 开启返回函数名和行号
	mLog.SetFormatter(&LogFormatter{})                  // 设置自己定义的formatter
	parseLevel, err := logrus.ParseLevel(config.CONFIG.Logger.Level)
	if err != nil {
		parseLevel = logrus.InfoLevel
	}
	mLog.SetLevel(parseLevel) // 设置level
	config.Log = mLog
}
