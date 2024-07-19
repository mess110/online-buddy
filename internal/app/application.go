package app

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type OnlineBuddy struct {
	Name      string
	Version   string
	Port      int
	DBAddress string

	Logger      *zap.Logger
	RedisClient *redis.Client
}

func Init() {
	// TODO: env variables
	dbAddress := "localhost:6379"
	port := 3000

	app := &OnlineBuddy{
		Name:      "online-buddy",
		Version:   "0.0.1",
		Port:      port,
		DBAddress: dbAddress,

		Logger:      NewLogger(),
		RedisClient: NewRedisClient(dbAddress),
	}

	app.Logger.Info(
		"Starting service",
		zap.String("name", app.Name),
		zap.String("version", app.Version),
	)

	serve(app)
}
