package app

import "go.uber.org/zap"

type OnlineBuddy struct {
	Logger  *zap.Logger
	Version string
	Port    int
	Name    string
}

func Init() {
	app := &OnlineBuddy{
		Logger:  NewLogger(),
		Version: "0.0.1",
		Port:    3000,
		Name:    "online-buddy",
	}

	app.Logger.Info("Starting service", zap.String("name", app.Name))

	serve(app)
}
