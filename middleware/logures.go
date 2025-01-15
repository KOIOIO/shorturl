package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logger 日志中间件
// todo 可考虑更换其他日志中间件
func Logger() gin.HandlerFunc {
	// 定义日志文件路径
	filePath := "log/log"
	//linkName := "latest_log.log"

	// 打开或创建日志文件
	scr, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	// 初始化日志记录器
	logger := logrus.New()

	// 设置日志输出目的地
	logger.Out = scr

	// 设置日志记录级别
	logger.SetLevel(logrus.DebugLevel)

	// 创建旋转日志写入器
	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
		//retalog.WithLinkName(linkName),
	)

	// 定义日志写入映射
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	// 创建新的日志钩子
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 添加钩子到日志记录器
	logger.AddHook(Hook)

	// 返回日志中间件处理函数
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		// 继续执行其他中间件或处理函数
		c.Next()
		// 计算请求处理时间
		stopTime := time.Since(startTime).Milliseconds()
		spendTime := fmt.Sprintf("%d ms", stopTime)
		// 获取主机名
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		// 获取HTTP响应状态码
		statusCode := c.Writer.Status()
		// 获取客户端IP地址
		clientIp := c.ClientIP()
		// 获取用户代理字符串
		userAgent := c.Request.UserAgent()
		// 获取响应数据大小
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		// 获取HTTP请求方法
		method := c.Request.Method
		// 获取请求路径
		path := c.Request.RequestURI

		// 创建日志条目并添加字段
		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		// 根据错误和状态码记录不同级别的日志
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
