package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisBloom struct {
	rdb *redis.Client
	key string
}

func NewRedisBloom(rdb *redis.Client, key string) *RedisBloom {
	return &RedisBloom{
		rdb: rdb,
		key: key,
	}
}

// 初始化布隆过滤器
func (rb *RedisBloom) Init() error {
	// BF.RESERVE 创建布隆过滤器
	// error_rate: 0.01 表示1%的错误率
	// initial_size: 1000000 预计元素数量
	return rb.rdb.Do(context.Background(), "BF.RESERVE", rb.key, 0.01, 1000000).Err()
}

// 添加元素
func (rb *RedisBloom) Add(url string) error {
	return rb.rdb.Do(context.Background(), "BF.ADD", rb.key, url).Err()
}

// 检查元素是否存在
func (rb *RedisBloom) MightContain(url string) (bool, error) {
	result, err := rb.rdb.Do(context.Background(), "BF.EXISTS", rb.key, url).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
