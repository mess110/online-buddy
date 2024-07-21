package app

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
)

type RedisDB struct {
	Client        *redis.Client
	ClusterClient *redis.ClusterClient

	ctx context.Context

	addresses []string
	isCluster bool
}

func NewRedisClient(logger *zap.Logger, addrs []string, isCluster bool) *RedisDB {
	var rdbc *redis.ClusterClient
	var rdb *redis.Client
	var pingErr error

	if isCluster {
		rdbc = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addrs,
		})

		_, pingErr = rdbc.Ping().Result()
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: addrs[0],
		})

		_, pingErr = rdb.Ping().Result()
	}

	if pingErr != nil {
		logger.Fatal("ping error", zap.Error(pingErr))
	}

	ctx := context.Background()

	return &RedisDB{
		Client:        rdb,
		ClusterClient: rdbc,

		ctx: ctx,

		addresses: addrs,
		isCluster: isCluster,
	}
}

func (db *RedisDB) Subscribe(channels ...string) *redis.PubSub {
	if db.isCluster {
		return db.ClusterClient.Subscribe(channels...)
	} else {
		return db.Client.Subscribe(channels...)
	}
}

func (db *RedisDB) Publish(channel string, message interface{}) *redis.IntCmd {
	if db.isCluster {
		return db.ClusterClient.Publish(channel, message)
	} else {
		return db.Client.Publish(channel, message)
	}
}

func (db *RedisDB) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second) // 10 minutes
	if db.isCluster {
		return db.ClusterClient.Set(key, value, exp).Err()
	} else {
		return db.Client.Set(key, value, exp).Err()
	}
}

func (db *RedisDB) Get(key string) (interface{}, error) {
	if db.isCluster {
		return db.ClusterClient.Get(key).Result()
	} else {
		return db.Client.Get(key).Result()
	}
}
