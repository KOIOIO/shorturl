package server

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	model "shorturl/modle"
	"shorturl/utils/errmsg"
	"time"
)

// Bloom 是一个全局的布隆过滤器实例
var Bloom = NewBloomFilter()

//func GenerateShortURLString(url string) string {
//        // 使用MD5哈希算法对URL进行哈希计算
//        hash := md5.Sum([]byte(url))
//        // 取哈希值的前8个字节作为短URL的字节切片
//        shortURLBytes := hash[:]
//        // 使用Base64编码将短URL字节切片转换为字符串
//        encoded := base64.URLEncoding.EncodeToString(shortURLBytes)
//        // 返回编码后的字符串的前8个字符作为短URL
//        return encoded[:8]
//}

func DeleteWithTime() {
	for {
		time.Sleep(time.Hour * 24)
		model.Db.Where("created_at < ?", time.Now().Add(-time.Hour*24)).Delete(&model.Shorturl{})
	}
}
func InDb(url string) bool {
	var shortURL model.Shorturl
	if err := model.Db.Where("url =?", url).First(&shortURL).Error; err == nil {
		// 已经生成过，直接返回短链
		return true
	}
	return false
}

func GenerateShortURL(url string, expiration string) (int, string) {
	if url == "" {
		// 如果URL为空，返回错误
		return errmsg.ERROR_URL_IS_NULL, ""
	}
	if InDb(url) {
		// 已经生成过，直接返回短链
		var shortURL model.Shorturl
		if err := model.Db.Where("url =?", url).First(&shortURL).Error; err == nil {
			// 已经生成过，直接返回短链
			return errmsg.SUCCESS, shortURL.Shorturl
		}
	}
	// 使用布隆过滤器进行短链生成检查
	if Bloom.MightContain(url) {
		// 可能已经生成过，进行精确检查
		var shortURL model.Shorturl
		if err := model.Db.Where("url =?", url).First(&shortURL).Error; err == nil {
			// 已经生成过，直接返回短链
			log.Println("已经生成过，直接返回短链")
			return errmsg.SUCCESS, shortURL.Shorturl
		}
	}
	hash := md5.Sum([]byte(url))
	shortURLBytes := hash[:]
	encoded := base64.URLEncoding.EncodeToString(shortURLBytes)

	// 生成新的短链
	//shortURLStr := GenerateShortURLString(url)
	shortURLStr := encoded[:8]

	// 解析过期时间
	var expireDuration time.Duration
	var err error
	if expiration != "" {
		expireDuration, err = time.ParseDuration(expiration)
		if err != nil {
			// 如果过期时间格式无效，返回错误
			return errmsg.ERROR_EXPIRATION_ID_WRONG, ""
		}
	}

	// 将短链信息存储到 MySQL
	shortURL := model.Shorturl{
		Shorturl: shortURLStr,
		Url:      url,
	}
	if err := model.Db.Create(&shortURL).Error; err != nil {
		// 如果存储到MySQL失败，返回错误

		return errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL, ""
	}

	// 将短链和原始URL存储到 Redis，并设置过期时间
	if err := model.Redis.Rdb.Set(model.Redis.Ctx, shortURLStr, url, expireDuration).Err(); err != nil {
		// 如果存储到Redis失败，返回错误

		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	// 将原始URL添加到布隆过滤器
	Bloom.Add(url)

	return errmsg.SUCCESS, shortURLStr
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
