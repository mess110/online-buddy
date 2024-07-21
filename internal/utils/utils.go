package utils

import (
	"encoding/json"
	"os"
	"strings"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return logger
}

func EnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func GetAddrs(logger *zap.Logger, value string) ([]string, bool) {
	isCluster := strings.Contains(value, "[")
	var addrs []string

	if isCluster {
		err := json.Unmarshal([]byte(value), &addrs)
		if err != nil {
			logger.Fatal("could not parse address", zap.Error(err))
		}
	} else {
		addrs = []string{value}
	}
	return addrs, isCluster

}
