package model

type BriefContent struct {
	ContentID int      `json:"contentID"`
	Title     string   `json:"title"`
	Cover     string   `json:"cover"`
	Time      int64    `json:"createTime"`
	ViewNum   int      `json:"viewNum"`
	User      MiniUser `json:"author"`
}

type DetailedContent struct {
	ContentID  int      `json:"contentID"`
	Title      string   `json:"title"`
	Time       int64    `json:"createTime"`
	User       MiniUser `json:"author"`
	Liked      bool     `json:"liked"`
	ViewNum    int      `json:"viewNum"`
	CommentNum int      `json:"commentNum"`
	LikeNum    int      `json:"likeNum"`
	Tags       []string `json:"tags"`
}
