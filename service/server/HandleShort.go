package server

import (
	"errors"

	"gorm.io/gorm"
	model2 "shorturl/modle"
	"shorturl/server/repository"
	"shorturl/utils/errmsg"
	"time"
)

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
