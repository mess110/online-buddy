package app

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/mess110/online-buddy/internal/datatypes"
	"go.uber.org/zap"
)

type WebsocketChannel struct {
	App        *OnlineBuddy
	Connection *websocket.Conn
	Channel    string
}

func NewWebsocketChannel(app *OnlineBuddy, connection *websocket.Conn, channel string) *WebsocketChannel {
	return &WebsocketChannel{
		App:        app,
		Connection: connection,
		Channel:    channel,
	}
}

func (wc *WebsocketChannel) handle() {
	allFriends := wc.App.FriendGraph.GetAllFriends(wc.Channel)
	sub := wc.App.RedisDB.Subscribe(allFriends...)
	defer sub.Close()
	ch := sub.Channel()

	message := datatypes.NewUserStatusMessage(wc.Channel, datatypes.OnlineStatus)
	wc.App.RedisDB.Set(wc.Channel, string(message.Status))

	err := wc.sendOnlineFriends()
	if err != nil {
		wc.App.Logger.Error("write json error", zap.Error(err))
		wc.disconnect()
		return
	}

	err = wc.publish(message)
	if err != nil {
		wc.App.Logger.Error("publish error", zap.Error(err))
		wc.disconnect()
		return
	}

	go wc.handleDisconnect()

	for msg := range ch {
		err := wc.sendUserStatus(msg)
		if err != nil {
			wc.App.Logger.Error("sendUserStatus error", zap.Error(err))
			wc.disconnect()
			return
		}
	}
}

func (wc *WebsocketChannel) handleDisconnect() {
	for {
		messageType, _, err := wc.Connection.ReadMessage()
		if err != nil {
			wc.App.Logger.Error("read error", zap.Int("message_type", messageType), zap.Error(err))
			wc.disconnect()
			return
		}
	}
}

func (wc *WebsocketChannel) sendUserStatus(msg *redis.Message) error {
	data := datatypes.UserStatus{}
	err := json.Unmarshal([]byte(msg.Payload), &data)
	if err != nil {
		return err
	}

	err = wc.Connection.WriteJSON(data)
	if err != nil {
		return err
	}
	return nil
}

func (wc *WebsocketChannel) sendOnlineFriends() error {
	onlineFriends := []string{}
	allFriends := wc.App.FriendGraph.GetAllFriends(wc.Channel)

	for i, friend := range allFriends {
		iface, err := wc.App.RedisDB.Get(friend)
		if err != nil && err != redis.Nil {
			wc.App.Logger.Error("all friends", zap.Error(err))
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

	onlineFriendsMessage := datatypes.NewOnlineFriendsMessage(wc.Channel, onlineFriends)
	err := wc.Connection.WriteJSON(onlineFriendsMessage)
	if err != nil {
		return err
	}
	return nil
}

func (wc *WebsocketChannel) publish(message any) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = wc.App.RedisDB.Publish(wc.Channel, messageJSON).Err()
	if err != nil {
		return err
	}
	return nil
}

func (wc *WebsocketChannel) disconnect() {
	channel := wc.Channel

	message := datatypes.NewUserStatusMessage(channel, datatypes.OfflineStatus)
	wc.App.RedisDB.Set(channel, string(message.Status))
	err := wc.publish(message)
	if err != nil {
		wc.App.Logger.Error("publish error", zap.Error(err))
	}
}
