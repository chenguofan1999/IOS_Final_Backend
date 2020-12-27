package model

type Comment struct {
	CommentID int       `json:"commentID"`
	ContentID int       `json:"contentID"`
	Text      string    `json:"text"`
	Time      int64     `json:"createTime"`
	LikeNum   int       `json:"likeNum"`
	ReplyNum  int       `json:"replyNum"`
	User      *MiniUser `json:"user"`
	Liked     bool      `json:"liked"`
}
