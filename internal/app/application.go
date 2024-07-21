package app

import (
	"go.uber.org/zap"

	"github.com/mess110/online-buddy/internal/datatypes"
	"github.com/mess110/online-buddy/internal/db"
	"github.com/mess110/online-buddy/internal/utils"
)

type AppConfig struct {
	Name    string
	Version string
	Host    string
	Port    string

	DBAddrs   []string
	IsCluster bool

	RedisDB     *db.RedisDB
	FriendGraph *datatypes.FriendGraph
}

var (
	logger = utils.NewLogger()
)

func Init() {
	dbAddress := utils.EnvOrDefault("REDIS_CONNECTION_URL", "localhost:6379")
	dbAddrs, isCluster := utils.GetAddrs(logger, dbAddress)

	/*
		Since its only a demo application, I am passing the RedisDB and
		FriendGraph using the config struct, instead there should be connection
		pools for both and they shouldn't be in the config object
	*/
	app := &AppConfig{
		Name:    "online-buddy",
		Version: "0.0.1",
		Host:    "0.0.0.0",
		Port:    utils.EnvOrDefault("API_PORT", "3000"),

		DBAddrs:   dbAddrs,
		IsCluster: isCluster,

		RedisDB:     db.NewRedisClient(logger, dbAddrs, isCluster),
		FriendGraph: datatypes.NewFriendGraph(),
	}

	logger.Info(
		"Starting service",
		zap.String("name", app.Name),
		zap.String("version", app.Version),
		zap.String("host", app.Host),
		zap.Strings("redis_db_addresses", app.DBAddrs),
		zap.Bool("redis_is_cluster", app.IsCluster),
	)

	serve(app)
}
