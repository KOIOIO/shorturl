package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"shorturl/utils"
)

// 创建一个Redis操作的结构体用于后续对redis的操作
type Rediscli struct {
	Ctx context.Context
	Rdb *redis.Client
}

// 创建一个全局变量，用于后续对redis的操作
var Redis *Rediscli

// 初始化redis
func InitRedis() {
	Redis = &Rediscli{}
	Redis.Ctx = context.Background()
	Redis.Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.RedisHost + utils.RedisPort,
		Password: utils.RedisPassword, // no password set
		DB:       0,                   // use default DB
	})

	pong, err := Redis.Rdb.Ping(Redis.Ctx).Result()
	if err != nil {
		panic(err)
		return
	}
	log.Printf("connect redis success, pong=%s\n", pong)
}
