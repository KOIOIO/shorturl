package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	model "shorturl/modle"
	"shorturl/server"
	"time"
)

func GenerateShortURL(c *gin.Context) {
	// 获取原始 URL
	url := c.PostForm("url")
	if url == "" {
		c.JSON(400, gin.H{"error": "url is required"})
		return
	}

	// 检查布隆过滤器
	if server.Bloom.MightContain(url) {
		// 可能已经生成过，进行精确检查
		var shortURL model.Shorturl
		if err := model.Db.Where("url =?", url).First(&shortURL).Error; err == nil {
			// 已经生成过，直接返回
			c.JSON(200, gin.H{"short_url": shortURL.Shorturl})
			return
		}
	}

	// 生成短链
	shortURLStr := server.GenerateShortURLString(url)
	expiration := c.PostForm("expiration") // 获取过期时间，例如 "1h", "30m", "1d"
	var expireDuration time.Duration
	var err error
	if expiration != "" {
		expireDuration, err = time.ParseDuration(expiration)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid expiration format"})
			return
		}
	}

	// 存储到 MySQL
	shortURL := model.Shorturl{
		Shorturl: shortURLStr,
		Url:      url,
	}
	if err := model.Db.Create(&shortURL).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to save to MySQL"})
		return
	}

	// 存储到 Redis 并设置过期时间
	if err := model.Redis.Rdb.Set(model.Redis.Ctx, shortURLStr, url, expireDuration).Err(); err != nil {
		c.JSON(500, gin.H{"error": "failed to save to Redis"})
		return
	}

	// 添加到布隆过滤器
	server.Bloom.Add(url)

	c.JSON(200, gin.H{"short_url": shortURLStr})
}

func HandleShortURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short URL is required"})
		return
	}
	// 从 Redis 中查找原始 URL
	originalURL, err := model.Redis.Rdb.Get(model.Redis.Ctx, shortURL).Result()
	if err == redis.Nil {
		// 短链不存在于 Redis 中，可能已过期或未生成
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found or expired"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch original URL"})
		return
	}
	// 重定向到原始 URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}
