package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	model "shorturl/modle"
	"shorturl/server"
	"shorturl/utils/errmsg"
	"time"
)

// GenerateShortURL 生成短链接口
// @Summary 生成短链
// @Description 通过POST请求生成短链
// @Tags URL Shortening
// @Accept x-www-form-urlencoded
// @Param url formData string true "原始URL"
// @Param expiration formData string false "过期时间"
// @Success 200 {object} gin.H{"short_url": "生成的短链"}
// @Failure 400 {object} gin.H{"error": "错误信息"}
// @Failure 500 {object} gin.H{"error": "错误信息"}
// @Router /shorten [post]
func GenerateShortURL(c *gin.Context) {
	// 获取原始 URL
	url := c.PostForm("url")
	// 获取过期时间，例如 "1h", "30m", "1d"
	expiration := c.PostForm("expiration")
	if url == "" {
		// 如果URL为空，返回错误
		c.JSON(400, gin.H{"error": "url is required"})
		return
	}

	// 使用布隆过滤器进行短链生成检查
	if server.Bloom.MightContain(url) {
		// 可能已经生成过，进行精确检查
		var shortURL model.Shorturl
		if err := model.Db.Where("url =?", url).First(&shortURL).Error; err == nil {
			// 已经生成过，直接返回短链
			c.JSON(200, gin.H{"short_url": shortURL.Shorturl})
			return
		}
	}

	// 生成新的短链
	shortURLStr := server.GenerateShortURLString(url)

	// 解析过期时间
	var expireDuration time.Duration
	var err error
	if expiration != "" {
		expireDuration, err = time.ParseDuration(expiration)
		if err != nil {
			// 如果过期时间格式无效，返回错误
			c.JSON(400, gin.H{"error": "时间出错,请检查参数"})
			return
		}
	}

	// 将短链信息存储到 MySQL
	shortURL := model.Shorturl{
		Shorturl: shortURLStr,
		Url:      url,
	}
	if err := model.Db.Create(&shortURL).Error; err != nil {
		// 如果存储到MySQL失败，返回错误
		c.JSON(500, gin.H{"error": "存储到数据库中失败"})
		return
	}

	// 将短链和原始URL存储到 Redis，并设置过期时间
	if err := model.Redis.Rdb.Set(model.Redis.Ctx, shortURLStr, url, expireDuration).Err(); err != nil {
		// 如果存储到Redis失败，返回错误
		c.JSON(500, gin.H{"error": "存储到缓存中失败"})
		return
	}

	// 将原始URL添加到布隆过滤器
	server.Bloom.Add(url)

	// 返回生成的短链
	c.JSON(200, gin.H{"short_url": shortURLStr})
}

// HandleShortURL 处理短链跳转
// @Summary 处理短链跳转
// @Description 通过GET请求处理短链跳转
// @Tags URL Shortening
// @Param shortURL path string true "短链"
// @Success 301 {object} string "原始URL"
// @Failure 400 {object} gin.H{"error": "错误信息"}
// @Router /s/{shortURL} [get]
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
