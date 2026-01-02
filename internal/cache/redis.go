package cache

import (
	"context"
	"time"
	"github.com/souls-syntax/Templates/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client				*redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client:client}
}


