package logic

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	limite_processer "example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/logic/limit_processer"

	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/models"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/logic/repository"
	"gorm.io/gorm"

	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateShortUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateShortUrlLogic {
	return &GenerateShortUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateShortUrlLogic) GenerateShortUrl(in *shortUrl.GenerateShortUrlRequest) (*shortUrl.GenerateShortUrlResponse, error) {
	ip := limite_processer.GetClientIP(l.ctx)
	if !limite_processer.AllowIP(ip) {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_RATE_LIMIT,
			Shortcode: "",
		}, errors.New("rate limit exceeded")
	}

	if in.Url == "" {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_URL_IS_NULL,
			Shortcode: "",
		}, errors.New("url is null")
	}

	// 先用布隆过滤器检查URL是否存在
	// 如果存在，则直接从MySQL中读取短链
	// 如果不存在，则生成新的短链ID，保存到MySQL和Redis中，并加入布隆过滤器
	exist, _ := l.svcCtx.RedisBloom.MightContain(in.Url)

	if exist {
		shortcode, err := l.ReadFormMysql(in.Url)
		if err == nil && shortcode != nil {
			return &shortUrl.GenerateShortUrlResponse{
				Code:      errmsg.SUCCESS,
				Shortcode: shortcode.Shorturl,
			}, nil
		}
	}

	// 生成短链ID
	id, err := repository.GetMyFlake().NextID()
	if err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR,
			Shortcode: "",
		}, errors.New("ID生成失败")
	}

	shortCode := repository.Base62Encode(id)

	// 解析过期时间
	var expireDuration time.Duration
	if in.Expiration != "" {
		// 使用自定义函数解析过期时间
		expireDuration, err = parseCustomDuration(in.Expiration)
		if err != nil {
			return &shortUrl.GenerateShortUrlResponse{
				Code:      errmsg.ERROR_EXPIRATION_ID_WRONG,
				Shortcode: "",
			}, errors.New("failed to parse expiration time")
		}
	}

	shorturl := models.Shorturl{
		ID:       id,
		Shorturl: shortCode,
		Url:      in.Url,
	}

	if err := l.svcCtx.DB.Create(&shorturl).Error; err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL,
			Shortcode: "",
		}, errors.New("fail to save to mysql")
	}

	if err := l.svcCtx.Redis.Rdb.Set(l.svcCtx.Redis.Ctx, shortCode, in.Url, expireDuration).Err(); err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_FAILED_SAVE_TO_REDIS,
			Shortcode: "",
		}, errors.New("fail to save to redis")
	}

	// 新增后再加入布隆过滤器
	_ = l.svcCtx.RedisBloom.Add(in.Url)

	return &shortUrl.GenerateShortUrlResponse{
		Code:      errmsg.SUCCESS,
		Shortcode: shortCode,
	}, nil
}

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

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func (l *GenerateShortUrlLogic) DeleteWithTime() error {
	err := l.svcCtx.DB.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&models.Shorturl{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *GenerateShortUrlLogic) ReadFormMysql(url string) (*models.Shorturl, error) {
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

func (l *GenerateShortUrlLogic) SaveToMysql(shorturl models.Shorturl) error {
	err := l.svcCtx.DB.Create(&shorturl).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
