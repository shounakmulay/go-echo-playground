package redis_service

import (
	"awesomeProject/configs"
	"awesomeProject/models/app_error"
	"awesomeProject/util/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func Get(key string) (string, error) {
	logger.Sugar.Infof("Getting value for key %s", key)
	var client = configs.RedisClient
	val, err := client.Get(context.Background(), key).Result()

	if err == redis.Nil {
		logger.Sugar.Infof("Key %s not found!", key)
		return "", app_error.RedisError{
			Message: "Key not found!",
		}
	}

	if err != nil {
		logger.Sugar.Infof("Error occured for key %s", key)
		return "", err
	}
	logger.Sugar.Infof("Cache hit for key %s, value = %s", key, val)
	return val, nil
}

func Set(key string, val string) {
	logger.Sugar.Infof("Setting value for key %s", key)
	var client = configs.RedisClient
	err := client.Set(context.Background(), key, val, 30*time.Second).Err()
	if err != nil {
		logger.Sugar.Infof("Error setting key %s", key)
		panic(err)
	}
	logger.Sugar.Infof("Value %s set for key %s in redis", val, key)
}
