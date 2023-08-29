package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"whisper/pkg/config"
)

var (
	RDB         *redis.Client
	KeyNotExist = "#KeyNotExist#"
)

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	})
}
