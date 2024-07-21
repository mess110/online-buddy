package app

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/mess110/online-buddy/internal/datatypes"
)

/*
Holds the connection for a channel and is responsible for subscribing and
publishing messages to the channel and corresponding connection.
*/
type WebsocketChannel struct {
	Config     *AppConfig
	Connection *websocket.Conn
	Channel    string
}

func NewWebsocketChannel(config *AppConfig, connection *websocket.Conn, channel string) *WebsocketChannel {
	return &WebsocketChannel{
		Config:     config,
		Connection: connection,
		Channel:    channel,
	}
}

func (wsc *WebsocketChannel) subscribe() {
	allFriends := wsc.Config.FriendGraph.GetAllFriends(wsc.Channel)
	sub := wsc.Config.RedisWriteDB.Subscribe(allFriends...)
	defer sub.Close()
	ch := sub.Channel()

	message := datatypes.NewUserStatusMessage(wsc.Channel, datatypes.OnlineStatus)
	wsc.Config.RedisWriteDB.Set(wsc.Channel, string(message.Status))

	err := wsc.sendOnlineFriends()
	if err != nil {
		logger.Error("write json error", zap.Error(err))
		wsc.disconnect()
		return
	}

	err = wsc.publish(message)
	if err != nil {
		logger.Error("publish error", zap.Error(err))
		wsc.disconnect()
		return
	}

	go wsc.handleDisconnect()

	for msg := range ch {
		err := wsc.sendUserStatus(msg)
		if err != nil {
			logger.Error("sendUserStatus error", zap.Error(err))
			wsc.disconnect()
			return
		}
	}
}

func (wsc *WebsocketChannel) publish(message any) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = wsc.Config.RedisWriteDB.Publish(wsc.Channel, messageJSON).Err()
	if err != nil {
		return err
	}
	return nil
}

func (wsc *WebsocketChannel) sendUserStatus(msg *redis.Message) error {
	data := datatypes.UserStatus{}
	err := json.Unmarshal([]byte(msg.Payload), &data)
	if err != nil {
		return err
	}

	err = wsc.Connection.WriteJSON(data)
	if err != nil {
		return err
	}
	return nil
}

func (wsc *WebsocketChannel) sendOnlineFriends() error {
	onlineFriends := []string{}
	allFriends := wsc.Config.FriendGraph.GetAllFriends(wsc.Channel)

	for i, friend := range allFriends {
		iface, err := wsc.Config.RedisReadDB.Get(friend)
		if err != nil && err != redis.Nil {
			logger.Error("all friends", zap.Error(err))
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

	onlineFriendsMessage := datatypes.NewOnlineFriendsMessage(wsc.Channel, onlineFriends)
	err := wsc.Connection.WriteJSON(onlineFriendsMessage)
	if err != nil {
		return err
	}
	return nil
}

func (wsc *WebsocketChannel) handleDisconnect() {
	for {
		messageType, _, err := wsc.Connection.ReadMessage()
		if err != nil {
			logger.Error("read error", zap.Int("message_type", messageType), zap.Error(err))
			wsc.disconnect()
			return
		}
	}
}

func (wsc *WebsocketChannel) disconnect() {
	channel := wsc.Channel

	message := datatypes.NewUserStatusMessage(channel, datatypes.OfflineStatus)
	wsc.Config.RedisWriteDB.Set(channel, string(message.Status))
	err := wsc.publish(message)
	if err != nil {
		logger.Error("publish error", zap.Error(err))
	}
}
