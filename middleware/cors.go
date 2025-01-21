package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	// 使用cors.New创建一个新的CORS中间件实例，并配置CORS设置
	return cors.New(
		cors.Config{
			//AllowAllOrigins:  true, // 如果启用，将允许所有域进行跨域请求
			AllowOrigins:     []string{"*"},                                                             // 等同于允许所有域名 #AllowAllOrigins:  true
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                       // 允许的HTTP方法
			AllowHeaders:     []string{"*", "Authorization"},                                            // 允许的HTTP头部
			ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type"}, // 暴露的HTTP头部
			AllowCredentials: true,                                                                      // 是否允许携带凭据的请求
			MaxAge:           12 * time.Hour,                                                            // 预取请求的结果在多少时间内有效
		},
	)
}
