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
		title VARCHAR,
		description VARCHAR,
		create_time BIGINT,
		cover_url VARCHAR,
		video_url VARCHAR,
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

// QueryBriefContentWithContentID 根据 contentID 生成 BriefContent 对象，返回错误如果 contentID 不存在
func QueryBriefContentWithContentID(contentID int) (*BriefContent, error) {
	if !CheckContentExist(contentID) {
		return nil, errors.New("no such content")
	}

	content := new(BriefContent)
	content.ContentID = contentID
	var userID int

	row := DB.QueryRow(`select title, cover_url, create_time, user_id
		from contents where content_id = ?`, contentID)
	// 已知 content 存在，Scan()不会返回错误
	err := row.Scan(&content.Title, &content.Cover, &content.Time, &userID)
	// TODO: 确认功能无误后请删除 panic 代码以及上面的 err
	if err != nil {
		panic(err)
	}

	// content 和 user 已知存在，不需要处理错误
	content.ViewNum, _ = QueryViewNumWithContentID(contentID)
	content.User, _ = QueryMiniUserWithUserID(userID)

	return content, nil
}
