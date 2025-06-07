package server

import (
	"errors"

	"gorm.io/gorm"
	model2 "shorturl/modle"
	"shorturl/server/repository"
	"shorturl/utils/errmsg"
	"time"
)

func GenerateShortURL(url string, expiration string) (int, string) {
	if url == "" {
		return errmsg.ERROR_URL_IS_NULL, ""
	}

	// 检查是否已存在
	if repository.Bloom.MightContain(url) {
		shortURL, err := repository.ReadFormMysql(url)
		if err == nil && shortURL != nil {
			return errmsg.SUCCESS, shortURL.Shorturl
		}
	}

	// 生成雪花ID
	id, err := repository.Flake.NextID()
	if err != nil {
		return errmsg.ERROR, ""
	}

	// 生成短码 (使用Base62编码)
	shortCode := repository.Base62Encode(id)

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

	repository.Bloom.Add(url)
	return errmsg.SUCCESS, shortCode
}

func HandleShort(shortCode string) (code int, OriginalURL string) {
	// 清理过期数据
	err := repository.DeleteWithTime()
	if err != nil {
		return errmsg.ERROR, ""
	}

	// 1. 先查Redis
	originalURL, err := repository.GetFormRedis(shortCode)
	if err == nil {
		return errmsg.SUCCESS, originalURL
	}

	// 2. Redis不存在，查数据库
	var shortURLRecord model2.Shorturl
	if _, err := repository.ReadFormMysql(shortCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR_NOT_FOUND_IN_MYSQL, ""
		}
		return errmsg.ERROR_OTHER_EMS, ""
	}

	// 3. 重新缓存到Redis (默认24小时)
	if err := repository.SaveToRedis(shortCode, shortURLRecord.Url, 24*time.Hour); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	return errmsg.SUCCESS, shortURLRecord.Url
}
