package repository

import (
	"errors"
	"gorm.io/gorm"
	model2 "shorturl/modle"
	"time"
)

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func DeleteWithTime() error {
	err := model2.Db.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&model2.Shorturl{}).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadFormMysql(url string) (*model2.Shorturl, error) {
	var shortURL model2.Shorturl
	err := model2.Db.Where("url = ?", url).First(&shortURL).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shortURL, nil
}

func SaveToMysql(shorturl model2.Shorturl) error {
	err := model2.Db.Create(&shorturl).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
