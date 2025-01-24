package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl/server"
	"shorturl/utils/errmsg"
)

// @Summary 生成短链接
// @Description 根据提供的原始URL和过期时间生成短链接
// @Tags ShortURL
// @Accept json
// @Produce json
// @Param url formData string true "原始URL"
// @Param expiration formData string true "过期时间，例如 '1h', '30m', '1d'"
// @Success 200
// @Failure 400
// @Router /generate [post]
func Generate(c *gin.Context) {
	// 获取原始 URL
	url := c.PostForm("url")
	// 获取过期时间，例如 "1h", "30m", "1d"
	expiration := c.PostForm("expiration")
	code, shortURLStr := server.GenerateShortURL(url, expiration)
	if code == errmsg.SUCCESS {
		c.JSON(200, gin.H{"short_url": shortURLStr})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "返回错误"})
		return
	}
}

// @Summary 处理短链接跳转
// @Description 根据提供的短链接参数重定向到原始URL
// @Tags ShortURL
// @Accept json
// @Produce json
// @Param shortURL path string true "短链接参数"
// @Success 301
// @Failure 400
// @Router /{shortURL} [get]
func HandleShortURL(c *gin.Context) {
	// 获取短链参数
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		// 如果短链为空，返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "提供的短链为空"})
		return
	}
	// 处理短链跳转
	code, originalURL := server.HandleShort(shortURL)
	if code != errmsg.SUCCESS {
		// 如果短链处理失败，返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "你要访问的短链找不到了"})
		return
	}
	// 重定向到原始URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}
