package controller

import (
	"ios/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetComments 查询评论集, 必要参数:
// 1. contentID: 获取某个 content 的所有评论
// 以下参数与上面的参数兼容
// 1. orderBy : likeNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
// 注意：要求登录状态下发送，即头部中应包含 authorization
func GetComments(c *gin.Context) {
	contentID, err := strconv.Atoi(c.Query("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "contentID (integer) query required",
		})
		return
	}

	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	orderBy := c.DefaultQuery("orderBy", "time")
	if orderBy == "likeNum" {
		orderBy = "like_num"
	} else if orderBy == "time" {
		orderBy = "comment_id"
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

	comments := model.QueryComments(loginUserID, contentID, orderBy, order)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   comments,
	})

}

type postCommentFormat struct {
	ContentID int    `json:"contentID" binding:"required"`
	Text      string `json:"text" binding:"required"`
}

// PostComment 发送评论，要求 header 中有 authorization 字段，请求体为 postCommentFormat 格式.
func PostComment(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var input postCommentFormat
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {contentID, text}",
		})
		return
	}

	if err := model.InsertComment(loginUserID, input.ContentID, input.Text); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteComment : 删除一条评论，该评论必须由自己发出
func DeleteComment(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "commentID (integer) required",
		})
		return
	}

	if err := model.DeleteCommentWithCommentID(loginUserID, commentID); err != nil {
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
