package server

import (
	"errors"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"

	model "shorturl/modle"
	"shorturl/utils/errmsg"
	"time"
)

// Bloom 是一个全局的布隆过滤器实例
var Bloom = NewBloomFilter()
var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte
	for num > 0 {
		remainder := num % 62
		num /= 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
	}
	return string(encoded)
}

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
	model.Db.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&model.Shorturl{})
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
		return errmsg.ERROR_URL_IS_NULL, ""
	}

	// 检查是否已存在
	if Bloom.MightContain(url) {
		var shortURL model.Shorturl
		if err := model.Db.Where("url = ?", url).First(&shortURL).Error; err == nil {
			return errmsg.SUCCESS, shortURL.Shorturl
		}
	}

	// 生成雪花ID
	id, err := flake.NextID()
	if err != nil {
		return errmsg.ERROR, ""
	}

	// 生成短码 (使用Base62编码)
	shortCode := base62Encode(id)

	// 解析过期时间
	var expireDuration time.Duration
	if expiration != "" {
		expireDuration, err = time.ParseDuration(expiration)
		if err != nil {
			return errmsg.ERROR_EXPIRATION_ID_WRONG, ""
		}
	}

	// 存储到数据库
	shortURL := model.Shorturl{
		ID:       id,        // 存储雪花ID
		Shorturl: shortCode, // 存储短码
		Url:      url,
	}
	if err := model.Db.Create(&shortURL).Error; err != nil {
		return errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL, ""
	}

	// 存储到Redis
	if err := model.Redis.Rdb.Set(model.Redis.Ctx, shortCode, url, expireDuration).Err(); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	Bloom.Add(url)
	return errmsg.SUCCESS, shortCode
}

// HandleShort 处理短URL，返回原始URL和状态码
// 该函数首先尝试从Redis中获取原始URL，如果失败则尝试从数据库中获取
// 如果在Redis和数据库中都找不到短URL，或者出现错误，将返回相应的错误码
func HandleShort(shortCode string) (code int, OriginalURL string) {
	// 清理过期数据
	DeleteWithTime()

	// 1. 先查Redis
	originalURL, err := model.Redis.Rdb.Get(model.Redis.Ctx, shortCode).Result()
	if err == nil {
		return errmsg.SUCCESS, originalURL
	}

	// 2. Redis不存在，查数据库
	var shortURLRecord model.Shorturl
	if err := model.Db.Where("shorturl = ?", shortCode).First(&shortURLRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR_NOT_FOUND_IN_MYSQL, ""
		}
		return errmsg.ERROR_OTHER_EMS, ""
	}

	// 3. 重新缓存到Redis (默认24小时)
	if err := model.Redis.Rdb.Set(model.Redis.Ctx, shortCode, shortURLRecord.Url, 24*time.Hour).Err(); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	return errmsg.SUCCESS, shortURLRecord.Url
}
