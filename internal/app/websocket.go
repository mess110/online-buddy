package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ctx      = context.Background()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	userStatus  = make(map[string]Presence)
	friendGraph = NewFriendGraph()
)

func HandleWebsocket(app *OnlineBuddy, w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error("error", zap.Error(err))
		dc(app, channel)
		return
	}
	defer conn.Close()

	friendsChannels := friendGraph[channel]
	sub := app.RedisClient.Subscribe(ctx, friendsChannels...)
	defer sub.Close()
	ch := sub.Channel()

	userStatus[channel] = OnlineStatus

	err = sendOnlineFriends(conn, channel)
	if err != nil {
		app.Logger.Error("write json error", zap.Error(err))
		dc(app, channel)
		return
	}

	message := NewUserStatusMessage(channel, OnlineStatus)
	err = publish(app, channel, message)
	if err != nil {
		app.Logger.Error("publish error", zap.Error(err))
		dc(app, channel)
		return
	}

	go handleDisconnect(app, conn, channel)

	for msg := range ch {
		err := sendUserStatus(conn, msg)
		if err != nil {
			app.Logger.Error("sendUserStatus error", zap.Error(err))
			dc(app, channel)
			return
		}
	}
}

func handleDisconnect(app *OnlineBuddy, conn *websocket.Conn, channel string) {
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			app.Logger.Error("read error", zap.Int("message_type", messageType), zap.Error(err))
			dc(app, channel)
			return
		}
	}
}

func sendUserStatus(conn *websocket.Conn, msg *redis.Message) error {
	data := UserStatus{}
	err := json.Unmarshal([]byte(msg.Payload), &data)
	if err != nil {
		return err
	}

	err = conn.WriteJSON(data)
	if err != nil {
		return err
	}
	return nil
}

func sendOnlineFriends(conn *websocket.Conn, channel string) error {
	friends := friendGraph[channel]
	onlineFriendsMessage := NewOnlineFriendsMessage(channel, friends)
	err := conn.WriteJSON(onlineFriendsMessage)
	if err != nil {
		return err
	}
	return nil
}

func publish(app *OnlineBuddy, channel string, message any) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = app.RedisClient.Publish(ctx, channel, messageJSON).Err()
	if err != nil {
		return err
	}
	return nil
}

func dc(app *OnlineBuddy, channel string) {
	userStatus[channel] = OfflineStatus
	message := NewUserStatusMessage(channel, OfflineStatus)
	err := publish(app, channel, message)
	if err != nil {
		app.Logger.Error("publish error", zap.Error(err))
	}
}