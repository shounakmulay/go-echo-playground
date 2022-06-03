package configs

import (
	"awesomeProject/util/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

func redisEnvUri() string {
	return os.Getenv("REDIS_URI")
}

func ConnectRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisEnvUri(),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	RedisClient = rdb
	logger.Sugar.Infof("Connected to Redis on %s", rdb.Options().Addr)
}

var RedisClient *redis.Client
