package model

import (
	"fmt"
)

// CreateContentTableIfNotExists Creates a Contents Table If Not Exists
func CreateContentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS contents(
		content_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		title VARCHAR,
		description VARCHAR,
		create_time BIGINT,
		cover_url VARCHAR,
		video_url VARCHAR,
		PRIMARY KEY (content_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create content table failed", err)
		return
	}
}

func CheckContentExist(contentID int) bool {
	var temp int
	row := DB.QueryRow("select content_id from contents where content_id = ?", contentID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}
