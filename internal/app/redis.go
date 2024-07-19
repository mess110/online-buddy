package app

import "github.com/go-redis/redis/v8"

func NewRedisClient(address string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})
	return rdb
}
