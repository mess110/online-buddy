package app

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisDB struct {
	Client  *redis.ClusterClient
	address string
	ctx     context.Context
}

func NewRedisClient(address string) *RedisDB {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: address,
	// })

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"master0:6079", "master1:6179", "master2:6279"},
	})

	ctx := context.Background()

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return &RedisDB{
		Client:  rdb,
		address: address,
		ctx:     ctx,
	}
}

func (db *RedisDB) Subscribe(channels ...string) *redis.PubSub {
	return db.Client.Subscribe(channels...)
}

func (db *RedisDB) Publish(channel string, message interface{}) *redis.IntCmd {
	return db.Client.Publish(channel, message)
}

func (db *RedisDB) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second) // 10 minutes
	return db.Client.Set(key, value, exp).Err()
}

func (db *RedisDB) Get(key string) (interface{}, error) {
	return db.Client.Get(key).Result()
}
