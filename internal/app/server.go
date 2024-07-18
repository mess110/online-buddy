package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func serve(app *OnlineBuddy) {
	port := app.Port

	http.HandleFunc("/ws/{channel}", func(w http.ResponseWriter, r *http.Request) {
		HandleWebsocket(app, w, r)
	})

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	app.Logger.Info("Listening", zap.Int("port", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		app.Logger.Error("could not listen", zap.Error(err))
	}
}
