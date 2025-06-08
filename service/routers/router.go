package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	v1 "shorturl/api/v1"
	_ "shorturl/docs"

	middleware2 "shorturl/middleware"
)

// @title           短链生成器
// @version         1.0
// @description     短链生成器

// @contact.name   Xiaoyu_Wang
// @contact.url    https://gitee.com/wang-wenyu-fdhfj/short-url
// @contact.email  2652777599@qq.com

// @host      localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func InitRouter() *gin.Engine {
	// 创建默认的gin路由器
	r := gin.Default()

	// 添加静态文件服务，将../web目录下的文件映射到/web路径
	r.Static("/web", "../web")

	// 使用中间件以处理跨域请求
	r.Use(middleware2.Cors())

	// 使用日志中间件记录请求信息
	r.Use(middleware2.Logger())

	r.POST("/generate", v1.Generate)
	r.GET("/:shortURL", v1.HandleShortURL)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务监听8080端口
	err := r.Run(":8080")

	// 检查启动过程中是否有错误发生
	if err != nil {
		panic(err)
	} else {
		println("服务启动成功")
	}
	return r
}
