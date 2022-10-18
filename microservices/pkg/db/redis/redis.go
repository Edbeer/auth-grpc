package redis

import (
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/Edbeer/microservices/internal/config"

)

func NewRedisClient(c *config.Config) *redis.Client {
	redisAddr := c.Redis.RedisAddr

	if redisAddr == "" {
		redisAddr = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		MinIdleConns: c.Redis.MinIdleConns,
		PoolSize:     c.Redis.PoolSize,
		PoolTimeout:  time.Duration(c.Redis.PoolTimeout) * time.Second,
		Password:     c.Redis.RedisPassword, // no password set
		DB:           c.Redis.DB,            // use default DB
	})

	return client
}