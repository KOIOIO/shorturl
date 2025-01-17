package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "shorturl/api/v1"
	"shorturl/middleware"
)

func InitRouter() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(middleware.Logger())
	r.POST("/generate", v1.GenerateShortURL)
	r.GET("/:shortURL", v1.HandleShortURL)
	r.Run(":8080")
}
