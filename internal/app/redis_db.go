package app

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisDB struct {
	ClusterClient *redis.ClusterClient

	ctx context.Context

	addresses []string
	isCluster bool
}

func NewRedisClient(addrs []string, isCluster bool) *RedisDB {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: address,
	// })

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})

	ctx := context.Background()

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return &RedisDB{
		ClusterClient: rdb,

		ctx: ctx,

		addresses: addrs,
		isCluster: isCluster,
	}
}

func (db *RedisDB) Subscribe(channels ...string) *redis.PubSub {
	return db.ClusterClient.Subscribe(channels...)
}

func (db *RedisDB) Publish(channel string, message interface{}) *redis.IntCmd {
	return db.ClusterClient.Publish(channel, message)
}

func (db *RedisDB) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second) // 10 minutes
	return db.ClusterClient.Set(key, value, exp).Err()
}

func (db *RedisDB) Get(key string) (interface{}, error) {
	return db.ClusterClient.Get(key).Result()
}
