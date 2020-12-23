package model

type MiniUser struct {
	UserID      int    `json:"userID"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar"`
	FollowerNum int    `json:"followerNum"`
}
