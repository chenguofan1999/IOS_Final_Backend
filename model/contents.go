package model

import (
	"errors"
	"fmt"
)

// CreateContentTableIfNotExists Creates a Contents Table If Not Exists
func CreateContentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS contents(
		content_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		title VARCHAR(64),
		description VARCHAR(255),
		create_time BIGINT,
		cover_url VARCHAR(255),
		video_url VARCHAR(255),
		PRIMARY KEY (content_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create content table failed", err)
		return
	}
}

func CheckContentExist(contentID int) bool {
	var temp int
	row := DB.QueryRow("select content_id from contents where content_id = ?", contentID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}

// InsertContent 插入一条 Content 记录，用户不存在或插入错误时返回错误
func InsertContent(title string, description string, coverURL string, videoURL string, time int64, userID int) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	_, err := DB.Exec(`insert into contents(user_id,title,description,create_time,cover_url,video_url)
		values (?,?,?,?,?,?)`, userID, title, description, time, coverURL, videoURL)
	if err != nil {
		fmt.Println(err)
		return errors.New("insert content failed")
	}

	return nil
}

// QueryBriefContentWithContentID 根据 contentID 生成 BriefContent 对象，返回 nil 如果 contentID 不存在
func QueryBriefContentWithContentID(contentID int) *BriefContent {
	if !CheckContentExist(contentID) {
		return nil
	}

	content := new(BriefContent)
	content.ContentID = contentID
	var userID int

	row := DB.QueryRow(`select title, cover_url, create_time, user_id
		from contents where content_id = ?`, contentID)
	// 已知 content 存在，Scan()不会返回错误
	err := row.Scan(&content.Title, &content.CoverURL, &content.Time, &userID)
	// TODO: 确认功能无误后请删除 panic 代码以及上面的 err
	if err != nil {
		panic(err)
	}

	// content 已知存在，不需要处理错误
	content.ViewNum, _ = QueryViewNumWithContentID(contentID)
	user := QueryMiniUserWithUserID(userID)
	if user != nil {
		content.User = user
	}

	return content
}

// QueryBriefContentsWithUserID 查询某个用户的所有内容(Brief),如果没有则返回空切片。不做错误处理。
func QueryBriefContentsWithUserID(userID int) []BriefContent {
	if !CheckUserExist(userID) {
		return []BriefContent{}
	}

	contents := make([]BriefContent, 0)
	rows, err := DB.Query(`select content_id from contents where user_id = ? order by content_id desc`, userID)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var contentID int
		rows.Scan(&contentID)

		content := QueryBriefContentWithContentID(contentID)
		if content != nil {
			contents = append(contents, *content)
		}
	}
	return contents
}

// QueryBriefContentsFollowing 查询某个用户关注的用户的所有内容(Brief),如果用户不存在或用户没有内容没有则返回空切片。不做错误处理。
func QueryBriefContentsFollowing(userID int) []BriefContent {
	if !CheckUserExist(userID) {
		return []BriefContent{}
	}

	contents := make([]BriefContent, 0)
	rows, err := DB.Query(`select content_id from contents, follow
		where followed_id = user_id and follower_id = ?`, userID)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var contentID int
		rows.Scan(&contentID)

		content := QueryBriefContentWithContentID(contentID)
		if content != nil {
			contents = append(contents, *content)
		}
	}
	return contents
}

// QueryBriefContentsPublic 获取所有公共内容
func QueryBriefContentsPublic() []BriefContent {
	contents := make([]BriefContent, 0)
	rows, err := DB.Query(`select content_id from contents`)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var contentID int
		rows.Scan(&contentID)

		content := QueryBriefContentWithContentID(contentID)
		if content != nil {
			contents = append(contents, *content)
		}
	}

	return contents
}

// QueryDetailedContent 用户 ID 为 currentUserID 的用户请求内容 ID 为 contentID 的内容
// currentUserID 用于获知是否已点赞, 以及记录浏览历史
// time用于记录浏览历史
func QueryDetailedContent(currentUserID int, contentID int, time int64) *DetailedContent {
	if !CheckUserExist(currentUserID) || !CheckContentExist(contentID) {
		return nil
	}

	// 已确定用户和内容均存在
	content := new(DetailedContent)
	content.ContentID = contentID

	row := DB.QueryRow(`select user_id, title, description, create_time, video_url 
		from contents where content_id = ?`, contentID)

	var userID int
	row.Scan(&userID, &content.Title, &content.Description, &content.Time, &content.VideoURL)

	content.User = QueryMiniUserWithUserID(userID)
	content.Liked, _ = QueryHasLikedContent(currentUserID, contentID)
	content.ViewNum, _ = QueryViewNumWithContentID(contentID)
	content.LikeNum, _ = QueryLikeNumWithContentID(contentID)
	content.CommentNum, _ = QueryCommentNumWithContentID(contentID)
	content.Tags, _ = QueryTagsWithContentID(contentID)

	// 假设获取内容详细信息总伴随着用户的查看内容，因此对此做一条记录
	InsertHistory(currentUserID, contentID, time)

	return content
}

// DeleteContentWithContentID 删除一条内容，返回错误如果该内容不存在
func DeleteContentWithContentID(contentID int) error {
	if !CheckContentExist(contentID) {
		return errors.New("no such content")
	}

	// 内容存在，因此不必检查 result
	_, err := DB.Exec(`delete from contents where content_id = ?`, contentID)
	if err != nil {
		return errors.New("delete content failed")
	}

	return nil
}
