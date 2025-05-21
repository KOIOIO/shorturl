package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"shorturl/utils"
)

// 创建一个Redis操作的结构体用于后续对redis的操作
type Rediscli struct {
	Ctx context.Context // 上下文，用于取消请求和传递请求级值
	Rdb *redis.Client   // Redis客户端，用于执行Redis命令
}

// 创建一个全局变量，用于后续对redis的操作
var Redis *Rediscli

// 初始化redis
// 该函数创建一个Redis客户端实例，并检查与Redis服务器的连接
func InitRedis() {
	Redis = &Rediscli{}
	Redis.Ctx = context.Background()
	Redis.Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.RedisConfig.RedisHost + utils.RedisConfig.RedisPort,
		Password: utils.RedisConfig.RedisPassword, // no password set
		DB:       0,                               // use default DB
	})

	pong, err := Redis.Rdb.Ping(Redis.Ctx).Result()
	if err != nil {
		panic(err)
		return
	}
	log.Printf("connect redis success, pong=%s\n", pong)
}
