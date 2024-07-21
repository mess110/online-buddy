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

	DBWriteAddrs   []string
	IsWriteCluster bool
	DBReadAddrs    []string
	IsReadCluster  bool

	RedisWriteDB *db.RedisDB
	RedisReadDB  *db.RedisDB
	FriendGraph  *datatypes.FriendGraph
}

var (
	logger = utils.NewLogger()
)

func Init() {
	dbWriteAddress := utils.EnvOrDefault("REDIS_WRITE_CONNECTION_URL", "localhost:6379")
	dbWriteAddrs, isWriteCluster := utils.GetAddrs(logger, dbWriteAddress)

	dbReadAddress := utils.EnvOrDefault("REDIS_READ_CONNECTION_URL", "localhost:6379")
	dbReadAddrs, isReadCluster := utils.GetAddrs(logger, dbReadAddress)

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

		DBWriteAddrs:   dbWriteAddrs,
		IsWriteCluster: isWriteCluster,
		DBReadAddrs:    dbReadAddrs,
		IsReadCluster:  isReadCluster,

		RedisWriteDB: db.NewRedisClient(logger, dbWriteAddrs, isWriteCluster),
		RedisReadDB:  db.NewRedisClient(logger, dbReadAddrs, isReadCluster),
		FriendGraph:  datatypes.NewFriendGraph(),
	}

	logger.Info(
		"Starting service",
		zap.String("name", app.Name),
		zap.String("version", app.Version),
		zap.String("host", app.Host),

		zap.Strings("redis_db_write_addresses", app.DBWriteAddrs),
		zap.Bool("redis_is_write_cluster", app.IsWriteCluster),
		zap.Strings("redis_db_read_addresses", app.DBReadAddrs),
		zap.Bool("redis_is_read_cluster", app.IsReadCluster),
	)

	serve(app)
}
