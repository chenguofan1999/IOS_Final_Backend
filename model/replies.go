package model

import (
	"errors"
	"fmt"
	"time"
)

func CheckReplyExist(replyID int) bool {
	var temp int
	row := DB.QueryRow("select reply_id from replies where reply_id = ?", replyID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}

// CreateReplyTableIfNotExists Creates a Reply Table If Not Exists
func CreateReplyTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS replies(
		reply_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		comment_id INT,
		reply_text VARCHAR(255),
		create_time BIGINT,
		PRIMARY KEY (reply_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create replies table failed", err)
		return
	}
}

// QueryReplyWithReplyID 根据 replyID 查询并构造 Reply 结构
func QueryReplyWithReplyID(currentUserID int, replyID int) *Reply {
	if !CheckReplyExist(replyID) {
		return nil
	}

	// reply 确认存在
	reply := new(Reply)
	reply.ReplyID = replyID
	var userID int

	row := DB.QueryRow(`select user_id, comment_id, reply_text, create_time
	from replies where reply_id = ?`, replyID)
	err := row.Scan(&userID, &reply.CommentID, &reply.Text, &reply.Time)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	reply.LikeNum, _ = QueryLikeNumWithReplyID(replyID)
	reply.Liked, _ = QueryHasLikedReply(currentUserID, replyID)
	user := QueryMiniUserWithUserID(userID)
	if user != nil {
		reply.User = user
	}

	return reply
}

// QueryReplies 查询一条 comment 的所有回复, 如果没有则返回空切片
func QueryReplies(currentUserID int, commentID int, orderBy string, order string) []Reply {
	if !CheckCommentExist(commentID) {
		return []Reply{}
	}

	replies := make([]Reply, 0)
	rows, err := DB.Query(`select reply_id from replies natural left outer join
		(select reply_id, count(1) as like_num
		from like_reply group by reply_id) as X
		where comment_id = ? order by `+orderBy+` `+order, commentID)

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var replyID int
		rows.Scan(&replyID)

		reply := QueryReplyWithReplyID(currentUserID, replyID)
		if reply != nil {
			replies = append(replies, *reply)
		}
	}
	return replies
}

// QueryReplyNumWithCommentID 查询评论的回复数，返回错误如果评论不存在
func QueryReplyNumWithCommentID(commentID int) (int, error) {
	if !CheckCommentExist(commentID) {
		return 0, errors.New("no such comment")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from replies where comment_id = ?) as X`, commentID)
	err := row.Scan(&num)

	// 如果没有记录, Scan() 会返回错误, 为正常情况
	if err != nil {
		return 0, nil
	}

	return num, nil
}

// InsertReply 插入一条回复，用户、评论不存在或插入错误时返回错误
func InsertReply(userID int, commentID int, text string) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	// 检查评论存在
	if !CheckCommentExist(commentID) {
		return errors.New("no such comment")
	}

	// 执行
	_, err := DB.Exec(`insert into replies(user_id, comment_id, reply_text, create_time)
	values(?,?,?,?)`, userID, commentID, text, time.Now().Unix())
	if err != nil {
		return errors.New("insert reply failed")
	}

	return nil
}

// DeleteReplyWithReplyID 删除一条回复，返回错误如果该回复不存在, 或用户无权限删除
func DeleteReplyWithReplyID(userID int, replyID int) error {
	if !CheckReplyExist(replyID) {
		return errors.New("no such reply")
	}

	// 回复存在，因此 0 row affected 代表评论的发出者不是此用户
	result, err := DB.Exec(`delete from replies where user_id = ? and reply_id = ?`, userID, replyID)
	if err != nil {
		return errors.New("delete reply failed")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no access")
	}

	return nil
}
