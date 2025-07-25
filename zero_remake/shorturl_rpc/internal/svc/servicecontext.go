package svc

import (
	"fmt"

	"example.com/shorturl/short-url/zero_remake/common/init_gorm"
	"example.com/shorturl/short-url/zero_remake/common/init_redis"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/config"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/logic/repository"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *init_redis.Rediscli
	LogPath    string
	RedisBloom *repository.RedisBloom
}

func NewServiceContext(c config.Config) *ServiceContext {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.DbUser,
		c.Mysql.DbPass,
		c.Mysql.DbHost,
		c.Mysql.DbPort,
		c.Mysql.DbName,
	)
	db := init_gorm.Init_gorm(dns)

	rds, _ := init_redis.InitRedis(c.BizRedis.RedisHost, c.BizRedis.RedisPort, c.BizRedis.RedisPass, 0)
	redisbloom := repository.NewRedisBloom(rds.Rdb, "shorturl:Bloom")
	return &ServiceContext{
		Config:     c,
		DB:         db,
		Redis:      rds,
		LogPath:    c.Log.Path,
		RedisBloom: redisbloom,
	}
}
