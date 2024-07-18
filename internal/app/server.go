package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func serve(app *OnlineBuddy) {
	port := app.Port

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	app.Logger.Info("Listening", zap.Int("port", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		app.Logger.Error("could not listen", zap.Error(err))
	}
}
