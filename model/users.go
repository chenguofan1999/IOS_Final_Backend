package model

import (
	"errors"
	"fmt"
)

// CreateUserTableIfNotExists Creates a Users Table If Not Exists
func CreateUserTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS users(
		user_id INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(32) UNIQUE,
		password VARCHAR(32),
		bio VARCHAR(64) DEFAULT '',
		avatar_url VARCHAR(128) DEFAULT '',
		PRIMARY KEY (user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed", err)
		return
	}
}

// InsertUser 插入一个用户，除了 username 和 password 以外都是默认值
func InsertUser(username string, password string) error {
	if username == "" || password == "" {
		return errors.New("Invalid string")
	}

	_, err := DB.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("Insert user failed,err:%v", err)
		return errors.New("User exists")
	}

	return nil
}

// QueryUserIDWithName 通过用户名查询用户 ID , error != nil 如果不存在
func QueryUserIDWithName(username string) (int, error) {
	row := DB.QueryRow("select user_id from users where username = ?", username)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, errors.New("no such user")
	}
	return userID, nil
}

// CheckUserExist 检查 followerID 用户存在
func CheckUserExist(userID int) bool {
	var temp int
	row := DB.QueryRow("select user_id from users where user_id = ?", userID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}

// UpdateBio 更新指定用户 ID 的用户的简介, 返回 err 如果用户不存在
func UpdateBio(userID int, newBio string) {
	DB.Exec("update users set bio = ? where user_id = ?", newBio, userID)
}

// UpdateAvatar 更新指定用户 ID 的用户的头像
func UpdateAvatar(userID int, newAvatarURL string) {
	DB.Exec("update users set avatar_url = ? where user_id = ?", newAvatarURL, userID)
}
