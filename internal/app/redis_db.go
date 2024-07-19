package app

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client  *redis.Client
	address string
	ctx     context.Context
}

func NewRedisClient(address string) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
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
	return db.Client.Subscribe(db.ctx, channels...)
}

func (db *RedisDB) Publish(channel string, message interface{}) *redis.IntCmd {
	return db.Client.Publish(db.ctx, channel, message)
}

func (db *RedisDB) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second) // 10 minutes
	return db.Client.Set(db.ctx, key, value, exp).Err()
}

func (db *RedisDB) Get(key string) (interface{}, error) {
	val, err := db.Client.Get(db.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (db *RedisDB) MGet(keys []string) ([]interface{}, error) {
	return db.Client.MGet(db.ctx, keys...).Result()
}
