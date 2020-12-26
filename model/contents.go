package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CreateContentTableIfNotExists Creates a Contents Table If Not Exists
func CreateContentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS contents(
		content_id INT NOT NULL AUTO_INCREMENT,
		user_id INT,
		title VARCHAR(64),
		description VARCHAR(255),
		create_time BIGINT,
		cover_url VARCHAR(255),
		video_url VARCHAR(255),
		PRIMARY KEY (content_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create content table failed", err)
		return
	}
}

// CheckContentExist 检查 contentID 标识的 content 是否存在
func CheckContentExist(contentID int) bool {
	var temp int
	row := DB.QueryRow("select content_id from contents where content_id = ?", contentID)
	err := row.Scan(&temp)
	if err != nil {
		return false
	}
	return true
}

// InsertContent 插入一条 Content 记录，用户不存在或插入错误时返回错误
func InsertContent(title string, description string, coverURL string, videoURL string, time int64, userID int) error {
	// 检查用户存在
	if !CheckUserExist(userID) {
		return errors.New("no such user")
	}

	_, err := DB.Exec(`insert into contents(user_id,title,description,create_time,cover_url,video_url)
		values (?,?,?,?,?,?)`, userID, title, description, time, coverURL, videoURL)
	if err != nil {
		fmt.Println(err)
		return errors.New("insert content failed")
	}

	return nil
}

// QueryBriefContentWithContentID 根据 contentID 生成 BriefContent 对象，返回 nil 如果 contentID 不存在
func QueryBriefContentWithContentID(contentID int) *BriefContent {
	if !CheckContentExist(contentID) {
		return nil
	}

	content := new(BriefContent)
	content.ContentID = contentID
	var userID int

	row := DB.QueryRow(`select title, cover_url, create_time, user_id
		from contents where content_id = ?`, contentID)
	// 已知 content 存在，Scan()不会返回错误
	err := row.Scan(&content.Title, &content.CoverURL, &content.Time, &userID)
	// TODO: 确认功能无误后请删除 panic 代码以及上面的 err
	if err != nil {
		panic(err)
	}

	// content 已知存在，不需要处理错误
	content.ViewNum, _ = QueryViewNumWithContentID(contentID)
	user := QueryMiniUserWithUserID(userID)
	if user != nil {
		content.User = user
	}

	return content
}

// QueryContents 是查询内容集合的统一接口
// mode: public / user / tag / followBy
// specifier: 当模式为 user / tag / followBy 时，specifier 分别表示用户名 / tag名 / 用户名
// orderBy: view_num / create_time
// order : asc / desc
// num: 条数
func QueryContents(mode string, specifier interface{}, orderBy string, order string, num int) []BriefContent {

	fmt.Println("Querying contents...")
	fmt.Println("mode: ", mode)
	fmt.Println("specifier: ", specifier)
	fmt.Println("orderBy: ", orderBy)
	fmt.Println("order: ", order)
	fmt.Println("num: ", num)

	// 创建视图 view_num
	DB.Exec(`create view view_num as
	select content_id, count(1) as view_num from contents join history using (content_id) group by content_id;`)

	var rows *sql.Rows
	var err error

	switch mode {
	case "public":
		rows, err = DB.Query(`select content_id from contents natural join view_num order by `+orderBy+` `+order+` limit ?`, num)
	case "user":
		rows, err = DB.Query(`select content_id from contents natural join view_num 
		where user_id = ? order by `+orderBy+` `+order+` limit ?`, specifier, num)
	case "tag":
		rows, err = DB.Query(`select content_id from content_tags natural join view_num 
		where tag_name = ? order by `+orderBy+` `+order+` limit ?`, specifier, num)
	case "follow":
		rows, err = DB.Query(`select content_id from contents natural join view_num
		join follow on (user_id = followed_id and follower_id = ?) order by `+orderBy+` `+order+` limit ?`, specifier, num)
	}

	contents := make([]BriefContent, 0)

	if err != nil {
		fmt.Println(err)
		return contents
	}

	for rows.Next() {
		var contentID int
		rows.Scan(&contentID)

		content := QueryBriefContentWithContentID(contentID)
		if content != nil {
			contents = append(contents, *content)
		}
	}

	// 撤销视图 view_num
	DB.Exec(`drop view view_num;`)

	return contents
}

// QueryDetailedContent 用户 ID 为 currentUserID 的用户请求内容 ID 为 contentID 的内容.
// 参数：
// 1. currentUserID 用于获知是否已点赞, 以及记录浏览历史.
func QueryDetailedContent(currentUserID int, contentID int) *DetailedContent {
	if !CheckUserExist(currentUserID) || !CheckContentExist(contentID) {
		return nil
	}

	// 已确定用户和内容均存在
	content := new(DetailedContent)
	content.ContentID = contentID

	row := DB.QueryRow(`select user_id, title, description, create_time, video_url 
		from contents where content_id = ?`, contentID)

	var userID int
	row.Scan(&userID, &content.Title, &content.Description, &content.Time, &content.VideoURL)

	content.User = QueryMiniUserWithUserID(userID)
	content.Liked, _ = QueryHasLikedContent(currentUserID, contentID)
	content.ViewNum, _ = QueryViewNumWithContentID(contentID)
	content.LikeNum, _ = QueryLikeNumWithContentID(contentID)
	content.CommentNum, _ = QueryCommentNumWithContentID(contentID)
	content.Tags, _ = QueryTagsWithContentID(contentID)

	// 假设获取内容详细信息总伴随着用户的查看内容，因此对此做一条记录
	InsertHistory(currentUserID, contentID, time.Now().Unix())

	return content
}

// DeleteContentWithContentID 删除一条内容，返回错误如果该内容不存在
func DeleteContentWithContentID(contentID int) error {
	if !CheckContentExist(contentID) {
		return errors.New("no such content")
	}

	// 内容存在，因此不必检查 result
	_, err := DB.Exec(`delete from contents where content_id = ?`, contentID)
	if err != nil {
		return errors.New("delete content failed")
	}

	return nil
}
