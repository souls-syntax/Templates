package handlers

import (

	"github.com/souls-syntax/Templates/internal/cache"
)

type ApiConfig struct {
	Cache *cache.RedisCache

