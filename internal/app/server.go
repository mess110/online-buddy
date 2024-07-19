package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func serve(app *OnlineBuddy) {
	host := app.Host
	port := app.Port

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		b, err := json.Marshal(map[string]string{"ping": "pong"})
		if err != nil {
			app.Logger.Error("error marshaling", zap.Error(err))
		}
		w.Write(b)
	})

	http.HandleFunc("/ws/{channel}", func(w http.ResponseWriter, r *http.Request) {
		HandleWebsocket(app, w, r)
	})

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	app.Logger.Info("Listening", zap.String("host", host), zap.String("port", port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		app.Logger.Error("could not listen", zap.Error(err))
	}
}
