package routers

import (
	"github.com/gin-gonic/gin"
	v1 "shorturl/api/v1"
	"shorturl/middleware"
)

//
//func InitRouter() {
//	r := gin.Default()
//	// 添加静态文件服务
//	r.Static("/web", "../web")
//	r.Use(middleware.Cors())
//	r.Use(middleware.Logger())
//	r.POST("/generate", v1.GenerateShortURL)
//	r.GET("/:shortURL", v1.HandleShortURL)
//
//	err := r.Run(":8080")
//	if err != nil {
//		panic(err)
//	} else {
//		log.Println("服务启动成功")
//	}
//}

// InitRouter 初始化路由配置
func InitRouter() {
	// 创建默认的gin路由器
	r := gin.Default()

	// 添加静态文件服务，将../web目录下的文件映射到/web路径
	r.Static("/web", "../web")

	// 使用中间件以处理跨域请求
	r.Use(middleware.Cors())

	// 使用日志中间件记录请求信息
	r.Use(middleware.Logger())

	// 配置POST请求路由，用于生成短URL
	r.POST("/generate", v1.GenerateShortURL)

	// 配置GET请求路由，用于处理短URL的访问
	r.GET("/:shortURL", v1.HandleShortURL)

	// 启动服务监听8080端口
	err := r.Run(":8080")

	// 检查启动过程中是否有错误发生
	if err != nil {
		panic(err)
	} else {
		println("服务启动成功")
	}
}
