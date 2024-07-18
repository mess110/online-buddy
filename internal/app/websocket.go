package app

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandleWebsocket(app *OnlineBuddy, w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error("error", zap.Error(err))
		return
	}
	defer conn.Close()

	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()
	ch := sub.Channel()

	friendGraph := NewFriendGraph()

	app.Logger.Info(
		"subscribed",
		zap.String("channel", channel),
		zap.Any("friends", friendGraph[channel]),
	)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				app.Logger.Error("read error", zap.Error(err))
				return
			}

			err = rdb.Publish(ctx, channel, string(msg)).Err()
			if err != nil {
				app.Logger.Error("public error", zap.Error(err))
				return
			}
		}
	}()

	for msg := range ch {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			app.Logger.Error("write error", zap.Error(err))
			return
		}
	}
}
