package model

import (
	"errors"
	"fmt"
)

func CreateContentTagsTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS content_tags(
		content_id INT,
		tag_name VARCHAR,
		PRIMARY KEY (content_id, tag_name),
		FOREIGN KEY (content_id) REFERENCES contents(content_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create content_tags table failed", err)
		return
	}
}

func InsertContentTag(contentID int, tagName string) error {
	// 检查 content 存在
	exist := CheckContentExist(contentID)
	if exist == false {
		return errors.New("no such content")
	}

	// 由于主键已经防止重复，不用检验 result, err 的唯一可能性是已有 tag
	_, err := DB.Exec("insert into content_tags(content_id,tag_name) values(?,?)", contentID, tagName)
	if err != nil {
		return errors.New("tag exists")
	}

	return nil
}

func DeleteContentTag(contentID int, tagName string) error {
	// 检查 content 存在
	exist := CheckContentExist(contentID)
	if exist == false {
		return errors.New("no such content")
	}

	// 没什么错误好发生的，用 result 检验是否本来不存在这样的 tag
	result, _ := DB.Exec("delete from content_tags where content_id = ? and tag_name = ?", contentID, tagName)
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("no such tag")
	}

	return nil
}

// QueryContentsWithTag 查询具有某个 tag 的所有内容，如果没有就为空
func QueryContentsWithTag(tagName string) []BriefContent {
	contents := make([]BriefContent, 0)

	rows, _ := DB.Query(`select content_id from contents,content_tags where tag_name = ?`, tagName)

	for rows.Next() {
		var contentID int
		rows.Scan(&contentID)

		content, _ := QueryBriefContentWithContentID(contentID)
		contents = append(contents, *content)
	}
	return contents
}
