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

func InitRouter() {
	r := gin.Default()
	// 添加静态文件服务
	r.Static("/web", "../web")
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	r.POST("/generate", v1.GenerateShortURL)
	r.GET("/:shortURL", v1.HandleShortURL)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	} else {
		println("服务启动成功")
	}
}
