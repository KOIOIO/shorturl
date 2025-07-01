package init_redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// 创建一个Redis操作的结构体用于后续对redis的操作
type Rediscli struct {
	Ctx context.Context // 上下文，用于取消请求和传递请求级值
	Rdb *redis.Client   // Redis客户端，用于执行Redis命令
}

// 创建一个全局变量，用于后续对redis的操作
//var Redis *Rediscli

// 初始化redis

func InitRedis(RedisHost string, RedisPort string, RedisPassword string, RedisDB int) (*Rediscli, error) {
	Redis := &Rediscli{
		Ctx: context.Background(),
		Rdb: redis.NewClient(&redis.Options{
			Addr:     RedisHost + ":" + RedisPort, // 添加:分隔符
			Password: RedisPassword,
			DB:       RedisDB,
		}),
	}

	// 检查连接
	if _, err := Redis.Rdb.Ping(Redis.Ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis连接失败: %v", err)
	}

	log.Printf("connect redis success")
	return Redis, nil
}
