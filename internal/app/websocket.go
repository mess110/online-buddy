package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/mess110/online-buddy/internal/datatypes"
	"go.uber.org/zap"
)

var (
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
		dc(app, channel)
		return
	}
	defer conn.Close()

	allFriends := app.FriendGraph.GetAllFriends(channel)
	sub := app.RedisDB.Subscribe(allFriends...)
	defer sub.Close()
	ch := sub.Channel()

	message := datatypes.NewUserStatusMessage(channel, datatypes.OnlineStatus)
	app.RedisDB.Set(channel, string(message.Status))

	err = sendOnlineFriends(app, conn, channel)
	if err != nil {
		app.Logger.Error("write json error", zap.Error(err))
		dc(app, channel)
		return
	}

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
	data := datatypes.UserStatus{}
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

func sendOnlineFriends(app *OnlineBuddy, conn *websocket.Conn, channel string) error {
	onlineFriends := []string{}
	allFriends := app.FriendGraph.GetAllFriends(channel)

	for i, friend := range allFriends {
		iface, err := app.RedisDB.Get(friend)
		if err != nil && err != redis.Nil {
			app.Logger.Error("all friends", zap.Error(err))
			return err
		}
		if iface != nil {
			val := iface.(string)
			if val == string(datatypes.OnlineStatus) {
				friend := allFriends[i]
				onlineFriends = append(onlineFriends, friend)
			}
		}
	}

	onlineFriendsMessage := datatypes.NewOnlineFriendsMessage(channel, onlineFriends)
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
	err = app.RedisDB.Publish(channel, messageJSON).Err()
	if err != nil {
		return err
	}
	return nil
}

func dc(app *OnlineBuddy, channel string) {
	message := datatypes.NewUserStatusMessage(channel, datatypes.OfflineStatus)
	app.RedisDB.Set(channel, string(message.Status))
	err := publish(app, channel, message)
	if err != nil {
		app.Logger.Error("publish error", zap.Error(err))
	}
}
