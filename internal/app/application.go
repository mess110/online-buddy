package app

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type OnlineBuddy struct {
	Name      string
	Version   string
	Host      string
	Port      string
	DBAddress string

	Logger      *zap.Logger
	RedisClient *redis.Client
}

func Init() {
	dbAddress := EnvOrDefault("REDIS_CONNECTION_URL", "localhost:6379")

	app := &OnlineBuddy{
		Name:      "online-buddy",
		Version:   "0.0.1",
		Host:      "0.0.0.0",
		Port:      EnvOrDefault("API_PORT", "3000"),
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
