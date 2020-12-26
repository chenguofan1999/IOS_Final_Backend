package controller

import (
	"errors"
	"ios/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserInfoByName ：在登录状态下，获取目标用户的详细信息。
func GetUserInfoByName(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	targetUserName := c.Param("username")
	targetUserID, err := model.QueryUserIDWithName(targetUserName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	// 已验证用户存在，此时 detailedInfo 不会为 nil
	detailedInfo := model.QueryDetailedUser(loginUserID, targetUserID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   *detailedInfo,
	})
}

func GetSelfInfo(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	detailedInfo := model.QueryDetailedUser(loginUserID, loginUserID)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   *detailedInfo,
	})
}

// GetUserIDByAuth ： 从请求报文的 Header 中解析出用户ID。发生错误时已向响应报文写入JSON,若返回错误请直接返回。
func GetUserIDByAuth(c *gin.Context) (int, error) {
	// 从请求头中取得认证字段
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "unauthorized",
		})
		return 0, errors.New("unauthorized")
	}

	// 获得用户名
	loginUserName := GetNameByToken(tokenString)

	// 获得用户ID
	loginUserID, err := model.QueryUserIDWithName(loginUserName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return 0, errors.New("no such user")
	}

	return loginUserID, nil
}
