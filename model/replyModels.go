package model

type Reply struct {
	ReplyID   int       `json:"replyID"`
	CommentID int       `json:"commentID"`
	Text      string    `json:"text"`
	Time      int64     `json:"createTime"`
	LikeNum   int       `json:"likeNum"`
	User      *MiniUser `json:"user"`
	Liked     bool      `json:"liked"`
}
