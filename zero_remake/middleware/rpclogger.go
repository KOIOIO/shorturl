package middleware

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// UnaryServerLogger 创建一元RPC调用的日志拦截器
func UnaryServerLogger(path string) grpc.UnaryServerInterceptor {
	logger := initLogger(path) // 复用之前创建的日志初始化函数

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		// 调用实际的处理方法
		resp, err = handler(ctx, req)

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录日志
		entry := logger.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration.String(),
		})

		if err != nil {
			entry.WithField("error", err.Error()).Error("RPC call failed")
		} else {
			entry.Info("RPC call succeeded")
		}

		return resp, err
	}
}

// StreamServerLogger 创建流式RPC调用的日志拦截器
func StreamServerLogger(path string) grpc.StreamServerInterceptor {
	logger := initLogger(path)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		err := handler(srv, ss)

		duration := time.Since(startTime)

		entry := logger.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration.String(),
		})

		if err != nil {
			entry.WithField("error", err.Error()).Error("RPC服务响应失败")
		} else {
			entry.Info("RPC服务响应成功")
		}

		return err
	}
}
