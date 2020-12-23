package model

import (
	"fmt"
)

// CreateReplyTableIfNotExists Creates a Reply Table If Not Exists
func CreateReplyTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS comments(
		reply_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		comment_id INT,
		reply_text VARCHAR,
		create_time BIGINT,
		PRIMARY KEY (reply_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create comment table failed", err)
		return
	}
	fmt.Println("Create comment table successed or it already exists")
}

func CheckReplyExist(replyID int) bool {
	var temp int
	row := DB.QueryRow("select reply_id from replys where reply_id = ?", replyID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}
