package app

type Presence string

const (
	OnlineStatus  Presence = "online"
	OfflineStatus Presence = "offline"
)

type UserStatus struct {
	UserID string   `json:"user_id"`
	Status Presence `json:"status"`
}

func NewUserStatusMessage(channel string, status Presence) *UserStatus {
	message := UserStatus{UserID: channel, Status: status}
	return &message
}

type UserFriendsOnline struct {
	UserID        string   `json:"user_id"`
	FriendsOnline []string `json:"friends_online"`
}

func NewOnlineFriendsMessage(channel string, friends []string) *UserFriendsOnline {
	message := UserFriendsOnline{UserID: channel, FriendsOnline: []string{}}

	for _, friend := range friends {
		if userStatus[friend] == OnlineStatus {
			message.FriendsOnline = append(message.FriendsOnline, friend)
		}
	}

	return &message
}
