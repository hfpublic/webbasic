package middlewares

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogWriter struct {
	Loger *rotatelogs.RotateLogs
}

var logWriter *LogWriter

func init() {

	logPath := "logs"
	logFile := "web.log"
	loger, err := rotatelogs.New(
		path.Join(logPath, logFile)+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(logFile),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		log.Fatalln("err", err)
	}

	logWriter = &LogWriter{
		Loger: loger,
	}
	log.SetOutput(logWriter.Loger)
}

// Logger 自定义log中间件
func Logger() gin.HandlerFunc {

	logClient := logrus.New()

	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln("err", err)
	}
	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter.Loger,
		logrus.DebugLevel: logWriter.Loger,
		logrus.FatalLevel: logWriter.Loger,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logClient.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
