package controller

import (
	"ios/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LikeReply(c *gin.Context) {
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

	if err := model.InsertLikeReply(loginUserID, replyID); err != nil {
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

func CancelLikeReply(c *gin.Context) {
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

	if err := model.DeleteLikeReply(loginUserID, replyID); err != nil {
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
