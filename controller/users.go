package controller

import (
	"errors"
	"fmt"
	"ios/model"
	"math/rand"
	"net/http"
	"path"

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

type bioInfo struct {
	Bio string `json:"bio" binding:"required"`
}

// UpdateUserBio : 更新当前用户的简介,请求体应为 JSON 形式,包含 bio 字段
func UpdateUserBio(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	// 获取请求体中的参数
	var info bioInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {bio}",
		})
		return
	}

	if err := model.UpdateBio(loginUserID, info.Bio); err != nil {
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

// UpdateUserAvatar : 更新当前用户的头像,请求体应为 Form-data 形式, 包含 key 为 avatar 的头像文件
func UpdateUserAvatar(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	// 读取文件
	avatarFile, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect Form-data: {avatar}",
		})
		return
	}

	// 检查后缀
	suffix := path.Ext(avatarFile.Filename)
	if suffix != ".jpg" && suffix != ".png" && suffix != ".jpeg" && suffix != ".svg" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "supported image type: jpg,png,jpeg,svg ",
		})
		return
	}

	// 检查图片大小不大于 8mb
	if avatarFile.Size > (8 << 20) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "image file too big (> 1mb)",
		})
		return
	}

	// 文件路径 和 URL
	randomNumberSuffix := rand.Intn(1000)
	filePath := fmt.Sprintf("/home/lighthouse/IOS_Files/avatars/user%d_avatar_%d%s", loginUserID, randomNumberSuffix, suffix)
	avatarURL := fmt.Sprintf("/static/avatars/user%d_avatar_%d%s", loginUserID, randomNumberSuffix, suffix)

	// 保存
	if err := c.SaveUploadedFile(avatarFile, filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "Save error",
		})
		return
	}

	// 更新数据库
	if err := model.UpdateAvatar(loginUserID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "update DB failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type tagFormat struct {
	Tag string `json:"tag" binding:"required"`
}

// AddTagForCurrentUser : 为当前用户增加一个 Tag, 请求体为 JSON 形式，包含 tag 字段
func AddTagForCurrentUser(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var newTag tagFormat
	if err := c.BindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {tag}",
		})
		return
	}

	if err := model.InsertUserTag(loginUserID, newTag.Tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// })
	GetTagsForCurrentUser(c)
}

// DeleteTagForCurrentUser : 为当前用户删除一个 Tag, 请求体为 JSON 形式，包含 tag 字段
func DeleteTagForCurrentUser(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var newTag tagFormat
	if err := c.BindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {tag}",
		})
		return
	}

	if err := model.DeleteUserTag(loginUserID, newTag.Tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// })
	GetTagsForCurrentUser(c)
}

// GetTagsForCurrentUser 获取已登录用户的 tags
func GetTagsForCurrentUser(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	tags, _ := model.QueryTagsWithUserID(loginUserID)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tags,
	})
}
