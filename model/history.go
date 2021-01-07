package model

import (
	"errors"
	"fmt"
)

func CreateHistoryTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS history(
		user_id INT,
		content_id INT,
		view_time BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create like_content table failed", err)
		return
	}
}

func InsertHistory(userID int, contentID int, time int64) error {
	// 确认用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	// 确认内容存在
	if !CheckContentExist(contentID) {
		return errors.New("no such content")
	}

	DB.Exec(`insert into history(user_id,content_id,view_time) values(?,?,?)`, userID, contentID, time)
	return nil
}

// QueryViewNumWithContentID 查询目标 content 的 view number，返回 err != nil 如果 content 不存在
func QueryViewNumWithContentID(contentID int) (int, error) {
	if !CheckContentExist(contentID) {
		return 0, errors.New("no such content")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from history where content_id = ?) as X`, contentID)
	err := row.Scan(&num)

	// 如果没有记录 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}
