package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-hclog"
)

type RedisCache struct {
	redisClient *redis.Client
	log         hclog.Logger
}

var ctx = context.Background()

func InitializeCacheClient() (*RedisCache, error) {

	rdsClient := redis.NewClient(&redis.Options{
		//Addr: "redis.default.svc.cluster.local:6379",
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdsClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to redis client - %s", err.Error())
	}
	log := hclog.Default()
	return &RedisCache{
		redisClient: rdsClient,
		log:         log,
	}, nil
}

func (c *RedisCache) Set(key string, value *Employee) error {

	c.log.Info("Set Redis Cache", "Key", key)

	json, err := json.Marshal(value)
	if err != nil {
		c.log.Error("Failed to Marshal cache data", "Key", key, "Error", err.Error())
		return err
	}

	err = c.redisClient.Set(ctx, key, json, 0).Err()
	if err != nil {
		c.log.Error("Set Redis Cache", "Key", key, "Error", err.Error())
		return err
	}
	return nil
}

func (c *RedisCache) Get(key string) (*Employee, error) {

	c.log.Info("Get Redis Cache", "Key", key)

	value, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		c.log.Error("Get Redis Cache", "Key", key, "Error", err.Error())
		return nil, err
	}

	employee := Employee{}
	err = json.Unmarshal([]byte(value), &employee)

	if err != nil {
		c.log.Error("Unable to Unmarshal Redis cache", "Key", key, "Error", err.Error())
		return nil, err
	}
	return &employee, nil
}

func (c *RedisCache) Del(key string) error {

	c.log.Info("Del Redis Cache", "Key", key)

	err := c.redisClient.Del(ctx, key).Err()
	if err != nil {
		c.log.Error("Unable to Del Redis Cache", "Key", key, "Error", err.Error())
		return err
	}

	return nil
}

func getKey(id int) string {
	return fmt.Sprintf("%s%d", "Key-", id)
}
