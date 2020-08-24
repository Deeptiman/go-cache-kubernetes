package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redisClient *redis.Client
}

var ctx = context.Background()

func InitializeCacheClient() (*RedisCache, error) {

	rdsClient := redis.NewClient(&redis.Options{
		Addr:     "redis.default.svc.cluster.local:6379", //"localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdsClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println(" <<<<  REDIS CLIENT NIL >>>> ")
		return nil, fmt.Errorf("Failed to connect to redis client - %s", err.Error())
	}

	return &RedisCache{
		redisClient: rdsClient,
	}, nil
}

func (cache *RedisCache) Set(key string, value *Employee) error {

	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Failed to Marshal cache data for Key %s - %s", key, err.Error())
	}

	err = cache.redisClient.Set(ctx, key, json, 0).Err()
	if err != nil {
		return fmt.Errorf("Failed to Set cache for Key %s - %s", key, err.Error())
	}
	return nil
}

func (cache *RedisCache) Get(key string) (*Employee, error) {

	value, err := cache.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to Get cache for Key %s - %s", key, err.Error())
	}

	employee := Employee{}
	err = json.Unmarshal([]byte(value), &employee)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode cached data for key %s , Error - %s", key, err.Error())
	}
	return &employee, nil
}
