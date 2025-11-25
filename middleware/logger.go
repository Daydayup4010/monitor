package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rota "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"path/filepath"
	"time"
	"uu/config"
)

func Logger() gin.HandlerFunc {
	logger := logrus.New()
	logDir := "log"
	fileName := "app.log"
	filePath := filepath.Join(logDir, fileName)
	linkName := filepath.Join(logDir, "last_log.log")
	scr, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0775)
	if err != nil {
		config.Log.Error(err)
	}
	logger.Out = scr

	parseLevel, err := logrus.ParseLevel(config.CONFIG.Logger.Level)
	if err != nil {
		parseLevel = logrus.InfoLevel
	}
	logger.SetLevel(parseLevel)
	logWriter, _ := rota.New(filePath+"%Y%m%d.log", rota.WithMaxAge(7*24*time.Hour), rota.WithRotationTime(24*time.Hour), rota.WithLinkName(linkName))
	wireMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.DebugLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.PanicLevel: logWriter,
		logrus.ErrorLevel: logWriter,
	}
	hook := lfshook.NewHook(wireMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(hook)
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0)))
		hostName := c.Request.Host
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.URL

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"SpendTime": spendTime,
			"Status":    statusCode,
			"ClientIP":  clientIP,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode < 500 && statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
