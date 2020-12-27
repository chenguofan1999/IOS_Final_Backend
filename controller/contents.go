package controller

import (
	"ios/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetContents 查询多条内容, 由 query 参数决定查询的模式:
// 这些参数是互斥的，即一条请求中只能 query 其中之一：
// 1. tag : 获取带有某个标签的内容 (tag = {tagName})
// 2. user : 获取指定用户的内容 (user = {userName})
// 3. follow : 当前用户关注的所有用户的内容 (follow = true)
// 4. self : 当前用户自己发的内容 (self = true)
// 5. 如果以上参数都没有，则为请求不经过筛选的公共内容
// 以下参数与上面的参数兼容
// 1. orderBy : viewNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
// 3. num : 指定条数, 默认 30
func GetContents(c *gin.Context) {
	tag := c.Query("tag")
	username := c.Query("user")
	follow := c.DefaultQuery("follow", "false")
	self := c.DefaultQuery("self", "false")

	orderBy := c.DefaultQuery("orderBy", "time")
	if orderBy == "viewNum" {
		orderBy = "view_num"
	} else if orderBy == "time" {
		orderBy = "content_id"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid query for orderBy",
		})
		return
	}

	order := c.DefaultQuery("order", "desc")
	if order != "desc" && order != "asc" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid query for order",
		})
		return
	}

	numStr := c.DefaultQuery("num", "30")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid query for num",
		})
		return
	}

	var contents []model.BriefContent
	// var err error

	if tag == "" && username == "" && follow == "false" && self == "false" {
		/* 公共内容 */
		contents = model.QueryContents("public", "_", orderBy, order, num)
	} else if tag != "" && username == "" && follow == "false" && self == "false" {
		/* 指定tag */
		contents = model.QueryContents("tag", tag, orderBy, order, num)
	} else if tag == "" && username != "" && follow == "false" && self == "false" {
		/* 指定user */
		userID, _ := model.QueryUserIDWithName(username)
		contents = model.QueryContents("user", userID, orderBy, order, num)
	} else if tag == "" && username == "" && follow == "true" && self == "false" {
		/* 我关注的 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("follow", loginUserID, orderBy, order, num)
	} else if tag == "" && username == "" && follow == "false" && self == "true" {
		/* 我的 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("user", loginUserID, orderBy, order, num)
	} else {
		/* 其他 */
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "only allow one of these query param at a time (tag / user / follow / self)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   contents,
	})
}

func GetContentByContentID(c *gin.Context) {
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expected contentID (integer)",
		})
		return
	}

	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	content := model.QueryDetailedContent(loginUserID, contentID)
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "content not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   content,
	})

}

// DeleteContent : 删除一条内容，该内容必须由自己发出
func DeleteContent(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "contentID (integer) required",
		})
		return
	}

	if err := model.DeleteContentWithContentID(loginUserID, contentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
