package server

import (
	"errors"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
	model2 "shorturl/modle"
	"shorturl/utils/errmsg"
	"time"
)

// Bloom 是一个全局的布隆过滤器实例
var Bloom = NewBloomFilter()

// flake 是一个全局的雪花ID生成器实例
var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

// base62Chars 定义了Base62编码使用的字符集
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// base62Encode 将一个整数转换为Base62编码的字符串
// 这个函数用于生成短URL的字符串表示
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

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func DeleteWithTime() {
	model2.Db.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&model2.Shorturl{})
}

// GenerateShortURL 生成短URL
// 参数:
//
//	url: 原始URL字符串
//	expiration: 过期时间字符串，格式如"24h"、"7d"等
//
// 返回值:
//
//	int: 错误码
//	string: 生成的短URL字符串
func GenerateShortURL(url string, expiration string) (int, string) {
	if url == "" {
		return errmsg.ERROR_URL_IS_NULL, ""
	}

	// 检查是否已存在
	if Bloom.MightContain(url) {
		var shortURL model2.Shorturl
		if err := model2.Db.Where("url = ?", url).First(&shortURL).Error; err == nil {
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
	shortURL := model2.Shorturl{
		ID:       id,        // 存储雪花ID
		Shorturl: shortCode, // 存储短码
		Url:      url,
	}
	if err := model2.Db.Create(&shortURL).Error; err != nil {
		return errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL, ""
	}

	// 存储到Redis
	if err := model2.Redis.Rdb.Set(model2.Redis.Ctx, shortCode, url, expireDuration).Err(); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	Bloom.Add(url)
	return errmsg.SUCCESS, shortCode
}

// HandleShort 处理短URL，返回原始URL和状态码
// 参数:
//
//	shortCode: 短URL的字符串表示
//
// 返回值:
//
//	int: 错误码
//	string: 原始URL字符串
//
// 该函数首先尝试从Redis中获取原始URL，如果失败则尝试从数据库中获取
// 如果在Redis和数据库中都找不到短URL，或者出现错误，将返回相应的错误码
func HandleShort(shortCode string) (code int, OriginalURL string) {
	// 清理过期数据
	DeleteWithTime()

	// 1. 先查Redis
	originalURL, err := model2.Redis.Rdb.Get(model2.Redis.Ctx, shortCode).Result()
	if err == nil {
		return errmsg.SUCCESS, originalURL
	}

	// 2. Redis不存在，查数据库
	var shortURLRecord model2.Shorturl
	if err := model2.Db.Where("shorturl = ?", shortCode).First(&shortURLRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR_NOT_FOUND_IN_MYSQL, ""
		}
		return errmsg.ERROR_OTHER_EMS, ""
	}

	// 3. 重新缓存到Redis (默认24小时)
	if err := model2.Redis.Rdb.Set(model2.Redis.Ctx, shortCode, shortURLRecord.Url, 24*time.Hour).Err(); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	return errmsg.SUCCESS, shortURLRecord.Url
}
