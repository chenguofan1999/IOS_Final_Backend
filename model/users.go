package model

import (
	"errors"
	"fmt"
)

// CreateUserTableIfNotExists Creates a Users Table If Not Exists
func CreateUserTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS users(
		user_id INT NOT NULL AUTO_INCREMENT,
		user_name VARCHAR(32) UNIQUE,
		password VARCHAR(32),
		bio VARCHAR(128) DEFAULT '',
		avatar_url VARCHAR(256) DEFAULT '',
		PRIMARY KEY (user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;`

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

	_, err := DB.Exec("insert INTO users(user_name,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("Insert user failed,err:%v", err)
		return errors.New("User exists")
	}

	return nil
}

// QueryUserIDWithName 通过用户名查询用户 ID , error != nil 如果不存在
func QueryUserIDWithName(username string) (int, error) {
	row := DB.QueryRow("select user_id from users where user_name = ?", username)
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

// QueryMiniUserWithUserID 根据用户 ID 获得 MiniUser 对象，具有足够用于辨别用户的属性，返回 nil 如果 user 不存在
func QueryMiniUserWithUserID(userID int) *MiniUser {
	if !CheckUserExist(userID) {
		return nil
	}

	// 已确定 user 存在,因此 QueryFollowerNumber 和 Scan() 不用处理错误
	user := new(MiniUser)
	user.UserID = userID
	user.FollowerNum, _ = QueryFollowerNumber(userID)

	row := DB.QueryRow(`select user_name, avatar_url from users where user_id = ?`, userID)
	row.Scan(&user.Username, &user.AvatarURL)

	return user
}
