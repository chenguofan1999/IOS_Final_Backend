package model

import (
	"fmt"
)

func CreateUserTagsTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS user_tags(
		user_id INT,
		tag_name VARCHAR,
		PRIMARY KEY (user_id, tag_name),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create user_tags table failed", err)
		return
	}
}
