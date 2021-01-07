package model

import (
	"errors"
	"fmt"
)

func CreateLikeReplyTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS like_reply(
		user_id INT,
		reply_id INT,
		PRIMARY KEY (user_id, reply_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (reply_id) REFERENCES replies(reply_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create like_reply table failed", err)
		return
	}
}

func QueryHasLikedReply(userID int, replyID int) (bool, error) {
	// 确认用户存在
	if !CheckUserExist(userID) {
		return false, errors.New("no such user")
	}

	// 确认回复存在
	if !CheckReplyExist(replyID) {
		return false, errors.New("no such reply")
	}

	// 查询 user 已 like reply
	var temp int
	row := DB.QueryRow("select 1 from like_reply where user_id = ? and reply_id = ?", userID, replyID)
	err := row.Scan(&temp)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func InsertLikeReply(userID int, replyID int) error {
	liked, err := QueryHasLikedReply(userID, replyID)
	if err != nil {
		return err
	} else if liked == true {
		return errors.New("already liked")
	}

	DB.Exec("insert into like_reply(user_id,reply_id) values(?,?)", userID, replyID)
	return nil
}

func DeleteLikeReply(userID int, replyID int) error {
	liked, err := QueryHasLikedReply(userID, replyID)
	if err != nil {
		return err
	} else if liked == false {
		return errors.New("did not like")
	}

	DB.Exec("delete from like_reply where user_id = ? and reply_id = ?", userID, replyID)
	return nil
}

func QueryLikeNumWithReplyID(replyID int) (int, error) {
	if !CheckReplyExist(replyID) {
		return 0, errors.New("no such reply")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from like_reply where reply_id = ?) as X`, replyID)
	err := row.Scan(&num)

	// 如果没有 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}
