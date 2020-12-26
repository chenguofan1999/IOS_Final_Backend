package controller

import (
	"ios/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFollowersByUserID : 获取指定用户的关注者列表
func GetFollowersByUserID(c *gin.Context) {
	username := c.Param("username")

	userID, err := model.QueryUserIDWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	followers, err := model.QueryFollowersWithUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   followers,
	})

}

// GetFollowingByUserID : 获取指定用户关注的用户的列表
func GetFollowingByUserID(c *gin.Context) {
	username := c.Param("username")

	userID, err := model.QueryUserIDWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	followers, err := model.QueryFollowingWithUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   followers,
	})

}

// FollowUser : 在登录状态下，当前用户关注目标用户
func FollowUser(c *gin.Context) {
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	username := c.Param("username")
	userID, err := model.QueryUserIDWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	err = model.InsertFollow(loginUserID, userID)
	if err != nil {
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

// UnfollowUser : 在登录状态下，当前用户取消关注目标用户
func UnfollowUser(c *gin.Context) {
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	username := c.Param("username")
	userID, err := model.QueryUserIDWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	err = model.DeleteFollow(loginUserID, userID)
	if err != nil {
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

// CheckFollowing : 在登录状态下，查询当前用户是否已经关注目标用户
func CheckFollowing(c *gin.Context) {
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	username := c.Param("username")
	userID, err := model.QueryUserIDWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	following, err := model.QueryHasFollowed(loginUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"following": following,
	})
}
