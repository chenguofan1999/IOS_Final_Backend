package model

import (
	"fmt"
)

func CreateContentTagsTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS user_tags(
		contentID INT,
		tag_name VARCHAR,
		PRIMARY KEY (contentID, tag_name),
		FOREIGN KEY (contentID) REFERENCES contents(contentID)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create content_tags table failed", err)
		return
	}
}
