package model

import (
	"errors"
	"fmt"
)

func CreateLikeContentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS like_content(
		user_id INT,
		content_id INT,
		PRIMARY KEY (user_id, content_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (content_id) REFERENCES contents(content_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create like_content table failed", err)
		return
	}
}

func QueryHasLikedContent(userID int, contentID int) (bool, error) {
	// 确认用户存在
	if !CheckUserExist(userID) {
		return false, errors.New("no such user")
	}

	// 确认内容存在
	if !CheckContentExist(contentID) {
		return false, errors.New("no such content")
	}

	// 查询 user 已 like content
	var temp int
	row := DB.QueryRow("select 1 from like_content where user_id = ? and content_id = ?", userID, contentID)
	err := row.Scan(&temp)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func InsertLikeContent(userID int, contentID int) error {
	liked, err := QueryHasLikedContent(userID, contentID)
	if err != nil {
		return err
	} else if liked == true {
		return errors.New("already liked")
	}

	DB.Exec("insert into like_content(user_id,content_id) values(?,?)", userID, contentID)
	return nil
}

func DeleteLikeContent(userID int, contentID int) error {
	liked, err := QueryHasLikedContent(userID, contentID)
	if err != nil {
		return err
	} else if liked == false {
		return errors.New("did not like")
	}

	DB.Exec("delete from like_content where user_id = ? and content_id = ?", userID, contentID)
	return nil
}

func QueryContentLikeNumber(contentID int) (int, error) {
	if !CheckContentExist(contentID) {
		return 0, errors.New("no such content")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from like_content where content_id = ?) as X`, contentID)
	err := row.Scan(&num)

	// 如果没有 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}
