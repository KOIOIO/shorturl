package svc

import (
	"example.com/shorturl/short-url/zero_remake/common/init_redis"
	"fmt"

	"example.com/shorturl/short-url/zero_remake/common/init_gorm"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *init_redis.Rediscli
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

	rds, _ := init_redis.InitRedis(c.BizRedis.RedisHost, c.BizRedis.RedisPort, c.BizRedis.RedisPass, 0.)

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,
	}
}
