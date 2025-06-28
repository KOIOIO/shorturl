package logic

import (
	"context"
	"errors"
	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/models"
	"gorm.io/gorm"
	"time"

	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleShortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleShortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleShortLogic {
	return &HandleShortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandleShortLogic) HandleShort(in *shortUrl.HandleShortRequest) (*shortUrl.HandleShortResponse, error) {
	err := l.DeleteWithTime()
	if err != nil {
		return &shortUrl.HandleShortResponse{
			Code:    errmsg.ERROR,
			LongUrl: "",
		}, errors.New("fail to delete expired records")
	}

	originalURL, err := l.GetFormRedis(in.Shortcode)
	if err == nil {
		return &shortUrl.HandleShortResponse{
			Code:    errmsg.SUCCESS,
			LongUrl: originalURL,
		}, nil
	}

	var shortURLRecord models.Shorturl
	if _, err := l.ReadFormMysql(in.Shortcode); err != nil {
		return &shortUrl.HandleShortResponse{
			Code:    errmsg.ERROR,
			LongUrl: "",
		}, errors.New("fail to get original URL")
	}

	if err := l.SaveToRedis(in.Shortcode, shortURLRecord.Url, time.Hour*24); err != nil {
		return &shortUrl.HandleShortResponse{
			Code:    errmsg.ERROR,
			LongUrl: "",
		}, errors.New("fail to save to redis")
	}

	return &shortUrl.HandleShortResponse{
		Code:    errmsg.SUCCESS,
		LongUrl: shortURLRecord.Url,
	}, nil
}

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func (l *HandleShortLogic) DeleteWithTime() error {
	err := l.svcCtx.DB.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&models.Shorturl{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *HandleShortLogic) ReadFormMysql(url string) (*models.Shorturl, error) {
	var shortURL models.Shorturl
	err := l.svcCtx.DB.Where("url = ?", url).First(&shortURL).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shortURL, nil
}

func (l *HandleShortLogic) SaveToMysql(shorturl models.Shorturl) error {
	err := l.svcCtx.DB.Create(&shorturl).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (l *HandleShortLogic) SaveToRedis(shortCode string, url string, expireDuration time.Duration) error {
	err := l.svcCtx.Redis.Rdb.Set(l.svcCtx.Redis.Ctx, shortCode, url, expireDuration).Err()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (l *HandleShortLogic) GetFormRedis(shortCode string) (string, error) {
	url, err := l.svcCtx.Redis.Rdb.Get(l.svcCtx.Redis.Ctx, shortCode).Result()
	if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
