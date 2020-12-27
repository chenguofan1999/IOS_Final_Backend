package controller

import (
	"ios/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetReplies 查询回复集, 必要参数:
// 1. commentID: 获取某个 comment 的所有回复
// 以下参数与上面的参数兼容
// 1. orderBy : likeNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
// 注意：要求登录状态下发送，即头部中应包含 authorization
func GetReplies(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Query("commentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "commentID (integer) query required",
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

	replies := model.QueryReplies(loginUserID, commentID, orderBy, order)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   replies,
	})
}

type postReplyFormat struct {
	CommentID int    `json:"commentID" binding:"required"`
	Text      string `json:"text" binding:"required"`
}

func PostReply(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var input postReplyFormat
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {commentID, text}",
		})
		return
	}

	if err := model.InsertReply(loginUserID, input.CommentID, input.Text); err != nil {
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

// DeleteReply : 删除一条回复，该回复必须由自己发出
func DeleteReply(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	replyID, err := strconv.Atoi(c.Param("replyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "replyID (integer) required",
		})
		return
	}

	if err := model.DeleteReplyWithReplyID(loginUserID, replyID); err != nil {
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
