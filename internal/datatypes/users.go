package datatypes

type Presence string

const (
	OnlineStatus  Presence = "online"
	OfflineStatus Presence = "offline"
)

type UserStatus struct {
	UserID string   `json:"user_id"`
	Status Presence `json:"status"`
}

type UserFriendsOnline struct {
	UserID        string   `json:"user_id"`
	FriendsOnline []string `json:"friends_online"`
}

func NewUserStatusMessage(channel string, status Presence) *UserStatus {
	message := UserStatus{UserID: channel, Status: status}
	return &message
}

func NewOnlineFriendsMessage(channel string, onlineFriends []string) *UserFriendsOnline {
	message := UserFriendsOnline{UserID: channel, FriendsOnline: onlineFriends}
	return &message
}
