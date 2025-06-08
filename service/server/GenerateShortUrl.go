package server

import (
	model2 "shorturl/modle"
	"shorturl/server/repository"
	"shorturl/utils/errmsg"
	"strconv"
	"strings"
	"time"
)

// parseCustomDuration 自定义函数，支持解析包含 'd' 单位的时间字符串
func parseCustomDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

func GenerateShortURL(url string, expiration string) (int, string) {
	if url == "" {
		return errmsg.ERROR_URL_IS_NULL, ""
	}

	if repository.Bloom.MightContain(url) {
		shortUrl, err := repository.ReadFormMysql(url)
		if err == nil && shortUrl != nil {
			return errmsg.SUCCESS, shortUrl.Shorturl
		}
	}

	id, err := repository.Flake.NextID()
	if err != nil {
		return errmsg.ERROR, ""
	}

	shortCode := repository.Base62Encode(id)

	// 解析过期时间
	var expireDuration time.Duration
	if expiration != "" {
		// 使用自定义函数解析过期时间
		expireDuration, err = parseCustomDuration(expiration)
		if err != nil {
			return errmsg.ERROR_EXPIRATION_ID_WRONG, ""
		}
	}

	shortUrl := model2.Shorturl{
		ID:       id,
		Shorturl: shortCode,
		Url:      url,
	}

	if err := model2.Db.Create(&shortUrl).Error; err != nil {
		return errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL, ""
	}

	if err := model2.Redis.Rdb.Set(model2.Redis.Ctx, shortCode, url, expireDuration).Err(); err != nil {
		return errmsg.ERROR_FAILED_SAVE_TO_REDIS, ""
	}

	repository.Bloom.Add(url)

	return errmsg.SUCCESS, shortCode
}
