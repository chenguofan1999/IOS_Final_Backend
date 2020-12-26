package model

import (
	"errors"
	"fmt"
)

// CreateFollowTableIfNotExists 构造 follow 表
func CreateFollowTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS follow(
		follower_id INT,
		followed_id INT,
		PRIMARY KEY (follower_id, followed_id),
		FOREIGN KEY (follower_id) REFERENCES users(user_id),
		FOREIGN KEY (followed_id) REFERENCES users(user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create Follow table failed", err)
		return
	}
}

// QueryHasFollowed 查询是否已经关注，返回错误如果任一用户不存在
func QueryHasFollowed(followerID int, followedID int) (bool, error) {
	// 检查用户存在
	if !CheckUserExist(followedID) || !CheckUserExist(followedID) {
		return false, errors.New("no such user")
	}

	// 查询已 follow
	var temp int
	row := DB.QueryRow("select 1 from follow where follower_id=? and followed_id=?", followerID, followedID)
	err := row.Scan(&temp)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// InsertFollow ：用户ID为 followerID 的用户 follow 用户ID为 followedID 的用户
func InsertFollow(followerID int, followedID int) error {
	// 不能 follow 自己
	if followerID == followedID {
		return errors.New("You can't follow yourself")
	}

	// 检查是否已经 follow
	following, err := QueryHasFollowed(followerID, followedID)
	if err != nil {
		return err
	} else if following == true {
		return errors.New("already following")
	}

	// 执行，执行语句不会出错
	DB.Exec("insert INTO follow(follower_id,followed_id) values(?,?)", followerID, followedID)

	fmt.Println(followerID, "follows", followedID)
	return nil
}

// DeleteFolloRelation ：用户ID为 followerID 的用户 unfollow 用户ID为 followedID 的用户
func DeleteFollow(followerID int, followedID int) error {
	// 不能 unfollow 自己
	if followerID == followedID {
		return errors.New("You can't unfollow yourself")
	}

	// 检查是否已经 follow
	following, err := QueryHasFollowed(followerID, followedID)
	if err != nil {
		return err
	} else if following == false {
		return errors.New("did not follow")
	}

	// 执行，执行语句不会出错
	DB.Exec("delete from follow where follower_id = ? and followed_id = ?", followerID, followedID)

	fmt.Println(followerID, "unfollows", followedID)
	return nil
}

// QueryFollowersWithUserID : 根据userID查询TA的关注者, 返回错误如果用户不存在或请求失败
func QueryFollowersWithUserID(userID int) ([]MiniUser, error) {
	if !CheckUserExist(userID) {
		return []MiniUser{}, errors.New("no such user")
	}

	followers := make([]MiniUser, 0)
	followerIDs, err := DB.Query(`select user_id,user_name,avatar_url 
		from users,follow where user_id = follower_id and followed_id = ?`, userID)

	if err != nil {
		return []MiniUser{}, errors.New("query follower failed")
	}

	for followerIDs.Next() {
		var user MiniUser
		followerIDs.Scan(&user.UserID, &user.Username, &user.AvatarURL)

		// 查询这个用户的关注者数
		num, _ := QueryFollowerNumWithUserID(user.UserID)
		user.FollowerNum = num

		followers = append(followers, user)
	}

	return followers, nil
}

// QueryFollowingWithUserID : 根据userID查询TA关注的人, 返回错误如果用户不存在或请求失败
func QueryFollowingWithUserID(userID int) ([]MiniUser, error) {
	if !CheckUserExist(userID) {
		return []MiniUser{}, errors.New("no such user")
	}

	following := make([]MiniUser, 0)
	followingIDs, err := DB.Query(`select user_id,user_name,avatar_url 
		from users,follow where user_id = followed_id and follower_id = ?`, userID)

	if err != nil {
		return []MiniUser{}, errors.New("query following failed")
	}

	for followingIDs.Next() {
		var user MiniUser
		followingIDs.Scan(&user.UserID, &user.Username, &user.AvatarURL)

		// 查询这个用户的关注者数
		num, _ := QueryFollowerNumWithUserID(user.UserID)
		user.FollowerNum = num

		following = append(following, user)
	}

	return following, nil
}

// QueryFollowerNumWithUserID 查询用户的关注者数目，返回 err != nil 如果用户不存在
func QueryFollowerNumWithUserID(userID int) (int, error) {
	if !CheckUserExist(userID) {
		return 0, errors.New("no such user")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from follow where followed_id = ?) as X`, userID)
	err := row.Scan(&num)

	// 如果没有 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}

// QueryFollowingNumWithUserID 查询关注的用户数目，返回 err != nil 如果用户不存在
func QueryFollowingNumWithUserID(userID int) (int, error) {
	if !CheckUserExist(userID) {
		return 0, errors.New("no such user")
	}

	var num int
	row := DB.QueryRow(`select count(1) from (select 1 from follow where follower_id = ?) as X`, userID)
	err := row.Scan(&num)

	// 如果没有 Scan() 会返回 err
	if err != nil {
		return 0, nil
	}

	return num, nil
}
