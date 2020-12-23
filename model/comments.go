package model

import (
	"fmt"
)

// CreateCommentTableIfNotExists Creates a Contents Table If Not Exists
func CreateCommentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS comments(
		comment_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		content_id INT,
		comment_text VARCHAR,
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
