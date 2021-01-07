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
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE
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

// InsertLikeContent : 插入一条喜欢内容的记录，返回错误如果用户/内容不存在，或已经喜欢
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

// DeleteLikeContent : 删除一条喜欢内容的记录，返回错误如果用户/内容不存在，或原本没有喜欢
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

// QueryLikeNumWithContentID 获取一条内容被赞的数目
func QueryLikeNumWithContentID(contentID int) (int, error) {
	if !CheckContentExist(contentID) {
		return 0, errors.New("no such content")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from like_content where content_id = ?) as X`, contentID)

	row.Scan(&num)
	return num, nil
}

// QueryLikeNumberWithUserID 获取用户(发布的内容)被赞的总数, 返回错误如果用户不存在
func QueryLikeNumberWithUserID(userID int) (int, error) {
	if !CheckUserExist(userID) {
		return 0, errors.New("no such user")
	}

	var num int
	row := DB.QueryRow(`select count(1) from  (select 1 
											   from like_content join contents using (content_id)
											   where contents.user_id = ?) as X`, userID)

	row.Scan(&num)
	return num, nil
}
