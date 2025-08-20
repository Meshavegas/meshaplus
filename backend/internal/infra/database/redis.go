package database

import (
	"backend/pkg/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// TODO: Add proper connection test
	return client, nil
}
