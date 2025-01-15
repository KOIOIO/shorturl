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

var Bloom = NewBloomFilter()

func GenerateShortURLString(url string) string {
	hash := md5.Sum([]byte(url))
	shortURLBytes := hash[:]
	encoded := base64.URLEncoding.EncodeToString(shortURLBytes)
	return encoded[:8]
}

//func HandleShort(shortURL string) (code int, OriginalURL string) {
//	// 从 Redis 中查找原始 URL
//	originalURL, err := model.Redis.Rdb.Get(model.Redis.Ctx, shortURL).Result()
//	if err == redis.Nil {
//		// 短链不存在于 Redis 中，可能已过期
//		return errmsg.ERROR, ""
//	} else if err != nil {
//		return errmsg.ERROR, ""
//	}
//	// 重定向到原始 URL
//	return errmsg.SUCCESS, originalURL
//}

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
				return errmsg.ERROR, ""
			}
			return errmsg.SUCCESS, originalURL
		} else if err == gorm.ErrRecordNotFound {
			// 数据库中也不存在该短链，返回错误
			return errmsg.ERROR, ""
		} else {
			// 数据库查询出现其他错误，返回错误
			return errmsg.ERROR, ""
		}
	} else if err != nil {
		// Redis 查找出现错误，返回错误
		return errmsg.ERROR, ""
	}
	// 重定向到原始 URL
	return errmsg.SUCCESS, originalURL
}
