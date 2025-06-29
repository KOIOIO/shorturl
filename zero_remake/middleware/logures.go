package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/rest"
)

// LoggerMiddleware 创建GoZero日志中间件
func LoggerMiddleware(path string) rest.Middleware {
	// 初始化日志配置
	logger := initLogger(path)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 记录请求开始时间
			startTime := time.Now()

			// 创建自定义ResponseWriter以捕获状态码和响应大小
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// 调用下一个处理器
			next(recorder, r)

			// 计算处理时间
			stopTime := time.Since(startTime).Milliseconds()
			spendTime := fmt.Sprintf("%d ms", stopTime)

			// 收集日志数据
			hostName, _ := os.Hostname()
			if hostName == "" {
				hostName = "unknown"
			}

			// 创建日志条目
			entry := logger.WithFields(logrus.Fields{
				"HostName":  hostName,
				"status":    recorder.statusCode,
				"SpendTime": spendTime,
				"Ip":        getClientIP(r),
				"Method":    r.Method,
				"Path":      r.URL.Path,
				"DataSize":  recorder.size,
				"Agent":     r.UserAgent(),
			})

			// 根据状态码记录不同级别日志
			switch {
			case recorder.statusCode >= 500:
				entry.Error("Server error")
			case recorder.statusCode >= 400:
				entry.Warn("Client error")
			default:
				entry.Info("Request processed")
			}
		}
	}
}

// 初始化日志记录器
func initLogger(path string) *logrus.Logger {
	filePath := path

	// 确保日志目录存在
	if err := os.MkdirAll("log", 0755); err != nil {
		fmt.Println("Failed to create log directory:", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 创建旋转日志
	logWriter, err := retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	// 设置日志输出
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return logger
}

// 获取客户端IP
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// 自定义ResponseWriter用于捕获状态码和响应大小
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}
