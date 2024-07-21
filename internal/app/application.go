package app

import (
	"go.uber.org/zap"
)

type OnlineBuddy struct {
	Name    string
	Version string
	Host    string
	Port    string

	DBAddrs   []string
	IsCluster bool

	Logger  *zap.Logger
	RedisDB *RedisDB

	FriendGraph *FriendGraph
}

func Init() {
	logger := NewLogger()

	dbAddress := EnvOrDefault("REDIS_CONNECTION_URL", "localhost:6379")
	dbAddrs, isCluster := GetAddrs(logger, dbAddress)

	app := &OnlineBuddy{
		Name:    "online-buddy",
		Version: "0.0.1",
		Host:    "0.0.0.0",
		Port:    EnvOrDefault("API_PORT", "3000"),

		DBAddrs:   dbAddrs,
		IsCluster: isCluster,

		Logger:  logger,
		RedisDB: NewRedisClient(logger, dbAddrs, isCluster),

		FriendGraph: NewFriendGraph(),
	}

	app.Logger.Info(
		"Starting service",
		zap.String("name", app.Name),
		zap.String("version", app.Version),
		zap.String("host", app.Host),
		zap.Strings("redis_db_addresses", dbAddrs),
		zap.Bool("redis_is_cluster", isCluster),
	)

	serve(app)
}
