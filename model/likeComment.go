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
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create like_comment table failed", err)
		return
	}
}

// QueryHasLikedComment 查询是否已经喜欢评论，返回错误如果用户不存在或评论不存在
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

// InsertLikeComment : 插入一条喜欢评论的记录，返回错误如果用户/评论不存在，或已经喜欢
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

// DeleteLikeComment : 删除一条喜欢评论的记录，返回错误如果用户/评论不存在，或原本没有喜欢
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

func QueryLikeNumWithCommentID(commentID int) (int, error) {
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
