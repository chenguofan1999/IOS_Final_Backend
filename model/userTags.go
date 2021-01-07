package model

import (
	"errors"
	"fmt"
)

func CreateUserTagsTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS user_tags(
		user_id INT,
		tag_name VARCHAR(32),
		PRIMARY KEY (user_id, tag_name),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create user_tags table failed", err)
		return
	}
}

// InsertUserTag 为一位用户增加 Tag, 返回错误如果用户不存在，或该用户已有该 Tag
func InsertUserTag(userID int, tagName string) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	// 由于主键已经防止重复，不用检验 result, err 的唯一可能性是已有 tag
	_, err := DB.Exec(`insert into user_tags(user_id, tag_name) values(?,?)`, userID, tagName)
	if err != nil {
		return errors.New("tag exists")
	}

	return nil
}

// DeleteUserTag 为一位用户删除 Tag, 返回错误如果用户不存在
func DeleteUserTag(userID int, tagName string) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	// 没什么错误好发生的，用 result 检验是否本来不存在这样的 tag
	result, _ := DB.Exec(`delete from user_tags where user_id = ? and tag_name = ?`, userID, tagName)
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("no such tag")
	}

	return nil
}

// QueryTagsWithUserID 根据用户 ID 查询其关注的 tag，返回错误如果用户不存在
func QueryTagsWithUserID(userID int) ([]string, error) {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return []string{}, errors.New("no such user")
	}

	tags := make([]string, 0)
	rows, _ := DB.Query(`select tag_name from user_tags where user_id = ?`, userID)

	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}

	return tags, nil
}
