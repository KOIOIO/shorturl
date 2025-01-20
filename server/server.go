package server

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	model "shorturl/modle"
	"shorturl/utils/errmsg"
	"time"
)

// Bloom 是一个全局的布隆过滤器实例
var Bloom = NewBloomFilter()

// GenerateShortURLString 将给定的长URL转换为短URL字符串
// 该函数使用MD5哈希和Base64编码来生成短URL
func GenerateShortURLString(url string) string {
	hash := md5.Sum([]byte(url))
	shortURLBytes := hash[:]
	encoded := base64.URLEncoding.EncodeToString(shortURLBytes)
	return encoded[:8]
}

// HandleShort 处理短URL，返回原始URL和状态码
// 该函数首先尝试从Redis中获取原始URL，如果失败则尝试从数据库中获取
// 如果在Redis和数据库中都找不到短URL，或者出现错误，将返回相应的错误码
func HandleShort(shortURL string) (code int, OriginalURL string) {
	// 从 Redis 中查找原始 URL
	originalURL, err := model.Redis.Rdb.Get(model.Redis.Ctx, shortURL).Result()
	if err == redis.Nil {
		// 短链不存在于 Redis 中，可能已过期，去数据库中查找
		var shortURLRecord model.Shorturl
		if err := model.Db.Where("shorturl =?", shortURL).First(&shortURLRecord).Error; err == nil {
			// 从数据库中找到了短链对应的原始 URL
			originalURL = shortURLRecord.Url
			// 将其重新添加到 Redis 中，假设设置过期时间为 1 小时，可根据需求修改
			err = model.Redis.Rdb.Set(model.Redis.Ctx, shortURL, originalURL, time.Hour).Err()
			if err != nil {
				// Redis 存储失败，返回错误
				return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
			}
			return errmsg.SUCCESS, originalURL
		} else if err == gorm.ErrRecordNotFound {
			// 数据库中也不存在该短链，返回错误
			return errmsg.ERROR_NOT_FOUND_IN_MYSQL, ""
		} else {
			// 数据库查询出现其他错误，返回错误
			return errmsg.ERROR_OTHER_EMS, ""
		}
	} else if err != nil {
		// Redis 查找出现错误，返回错误
		return errmsg.ERROR, ""
	}
	// 重定向到原始 URL
	return errmsg.SUCCESS, originalURL
}
