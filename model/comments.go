package model

import (
	"errors"
	"fmt"
	"time"
)

// CreateCommentTableIfNotExists Creates a Contents Table If Not Exists
func CreateCommentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS comments(
		comment_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		content_id INT,
		comment_text VARCHAR(255),
		create_time BIGINT,
		PRIMARY KEY (comment_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE,
		FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create comment table failed", err)
		return
	}
}

func CheckCommentExist(commentID int) bool {
	var temp int
	row := DB.QueryRow("select comment_id from comments where comment_id = ?", commentID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}

// InsertComment 插入一条评论，用户、内容不存在或插入错误时返回错误
func InsertComment(userID int, contentID int, text string) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	// 检查内容存在
	if !CheckContentExist(contentID) {
		return errors.New("no such content")
	}

	// 执行
	_, err := DB.Exec(`insert into comments(user_id,content_id,comment_text,create_time)
		values(?,?,?,?)`, userID, contentID, text, time.Now().Unix())
	if err != nil {
		return errors.New("insert comment failed")
	}

	return nil
}

// QueryCommentWithCommentID 根据 commentID 查询并构造 Comment 结构
func QueryCommentWithCommentID(currentUserID int, commentID int) *Comment {
	if !CheckCommentExist(commentID) {
		return nil
	}

	comment := new(Comment)
	comment.CommentID = commentID
	var userID int

	row := DB.QueryRow(`select user_id,content_id,comment_text,create_time 
	from comments where comment_id = ?`, commentID)
	err := row.Scan(&userID, &comment.ContentID, &comment.Text, &comment.Time)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	comment.LikeNum, _ = QueryLikeNumWithCommentID(commentID)
	comment.ReplyNum, _ = QueryReplyNumWithCommentID(commentID)
	comment.Liked, _ = QueryHasLikedComment(currentUserID, commentID)

	user := QueryMiniUserWithUserID(userID)
	if user != nil {
		comment.User = user
	}

	return comment
}

// QueryCommentNumWithContentID 查询一条内容的评论数，返回 err != nil 当contentID不存在
func QueryCommentNumWithContentID(contentID int) (int, error) {
	if !CheckContentExist(contentID) {
		return 0, errors.New("no such content")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from comments where content_id = ?) as X`, contentID)
	err := row.Scan(&num)

	// 如果没有记录 Scan() 会返回 err ，属于正常情况
	if err != nil {
		return 0, nil
	}

	return num, nil
}

// QueryComments 根据 contentID 查询对应内容的所有评论. 如果 content 不存在或没有评论则返回空切片.
func QueryComments(currentUserID int, contentID int, orderBy string, order string) []Comment {
	if !CheckContentExist(contentID) {
		return []Comment{}
	}

	comments := make([]Comment, 0)
	rows, err := DB.Query(`select comment_id from comments natural left outer join 
		(select comment_id, count(1) as like_num 
		from like_comment group by comment_id) as X 
		where content_id = ? order by `+orderBy+` `+order, contentID)
	if err != nil {
		fmt.Println(err)
		return comments
	}

	for rows.Next() {
		var commentID int
		rows.Scan(&commentID)

		comment := QueryCommentWithCommentID(currentUserID, commentID)
		if comment != nil {
			comments = append(comments, *comment)
		}
	}
	return comments
}

// DeleteCommentWithCommentID 删除一条评论，返回错误如果该评论不存在, 或用户无权限删除
func DeleteCommentWithCommentID(userID int, commentID int) error {
	if !CheckCommentExist(commentID) {
		return errors.New("no such comment")
	}

	// 评论存在，因此 0 row affected 代表评论的发出者不是此用户
	result, err := DB.Exec(`delete from comments where user_id = ? and comment_id = ?`, userID, commentID)
	if err != nil {
		return errors.New("delete comment failed")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no access")
	}

	return nil
}
