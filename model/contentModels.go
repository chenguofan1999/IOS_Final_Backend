package model

type BriefContent struct {
	ContentID int       `json:"contentID"`
	Title     string    `json:"title"`
	Duration  int       `json:"duration"`
	CoverURL  string    `json:"cover"`
	Time      int64     `json:"createTime"`
	ViewNum   int       `json:"viewNum"`
	User      *MiniUser `json:"user"`
}

type DetailedContent struct {
	ContentID   int       `json:"contentID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Time        int64     `json:"createTime"`
	VideoURL    string    `json:"video"`
	User        *MiniUser `json:"user"`
	Liked       bool      `json:"liked"`
	ViewNum     int       `json:"viewNum"`
	CommentNum  int       `json:"commentNum"`
	LikeNum     int       `json:"likeNum"`
	Tags        []string  `json:"tags"`
}
