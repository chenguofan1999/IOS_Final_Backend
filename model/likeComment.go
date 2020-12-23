package model

import (
	"errors"
	"fmt"
)

func CreateLikeCommentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS like_comment(
		user_id INT,
		comment_id INT,
		PRIMARY KEY (user_id, comment_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create like_comment table failed", err)
		return
	}
}

func QueryHasLikedComment(userID int, commentID int) (bool, error) {
	// 确认用户存在
	if !CheckUserExist(userID) {
		return false, errors.New("no such user")
	}

	// 确认评论存在
	if !CheckCommentExist(commentID) {
		return false, errors.New("no such comment")
	}

	// 查询 user 已 like comment
	var temp int
	row := DB.QueryRow("select 1 from like_comment where user_id = ? and comment_id = ?", userID, commentID)
	err := row.Scan(&temp)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func InsertLikeComment(userID int, commentID int) error {
	liked, err := QueryHasLikedComment(userID, commentID)
	if err != nil {
		return err
	} else if liked == true {
		return errors.New("already liked")
	}

	DB.Exec("insert into like_comment(user_id,comment_id) values(?,?)", userID, commentID)
	return nil
}

func DeleteLikeComment(userID int, commentID int) error {
	liked, err := QueryHasLikedComment(userID, commentID)
	if err != nil {
		return err
	} else if liked == false {
		return errors.New("did not like")
	}

	DB.Exec("delete from like_comment where user_id = ? and comment_id = ?", userID, commentID)
	return nil
}

func QueryCommentLikeNumber(commentID int) (int, error) {
	if !CheckCommentExist(commentID) {
		return 0, errors.New("no such comment")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from like_comment where comment_id = ?) as X`, commentID)
	err := row.Scan(&num)

	// 如果没有 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}
