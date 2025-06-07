package repository

import (
	model2 "shorturl/modle"
	"time"
)

func SaveToRedis(shortCode string, url string, expireDuration time.Duration) error {
	err := model2.Redis.Rdb.Set(model2.Redis.Ctx, shortCode, url, expireDuration).Err()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetFormRedis(shortCode string) (string, error) {
	url, err := model2.Redis.Rdb.Get(model2.Redis.Ctx, shortCode).Result()
	if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
