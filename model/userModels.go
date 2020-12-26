package model

type MiniUser struct {
	UserID      int    `json:"userID"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar"`
	FollowerNum int    `json:"followerNum"`
}

type DetailedUser struct {
	UserID       int    `json:"userID"`
	Bio          string `json:"bio"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatar"`
	FollowerNum  int    `json:"followerNum"`
	FollowingNum int    `json:"followingNum"`
	LikeNum      int    `json:"likeNum"`
	FollowedByMe bool   `json:"followedByMe"`
}
