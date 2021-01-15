package controller

import (
	"fmt"
	"ios/model"
	"math/rand"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetContents 查询多条内容, 由 query 参数决定查询的模式:
// 这些参数是互斥的，即一条请求中只能 query 其中之一：
// 1. tag : 获取带有某个标签的内容 (tag = {tagName})
// 2. user : 获取指定用户的内容 (user = {userName})
// 3. search : 搜索标题含特定字符串的内容 (tag = {searchStr})
// 4. follow : 当前用户关注的所有用户的内容 (follow = true)
// 5. self : 当前用户自己发的内容 (self = true)
// 6. history : 获取自己的观看记录 (history = true)
// 7. allTags : 获取自己关注的全部tag的内容 (allTags = true)
// 8. likedBy : 获取某人喜欢的全部内容 (likedBy = {username})
// 9. 如果以上参数都没有，则为请求不经过筛选的公共内容
// 以下参数与上面的参数兼容
// 1. orderBy : viewNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
// 3. num : 指定条数, 默认 30
func GetContents(c *gin.Context) {

	// count 用于计数互斥的参数
	count := 0

	tag := c.Query("tag")
	if tag != "" {
		count++
	}
	username := c.Query("user")
	if username != "" {
		count++
	}
	search := c.Query("search")
	if search != "" {
		count++
	}
	likedBy := c.Query("likedBy")
	if likedBy != "" {
		count++
	}
	follow := c.DefaultQuery("follow", "false")
	if follow == "true" {
		count++
	}
	self := c.DefaultQuery("self", "false")
	if self == "true" {
		count++
	}
	history := c.DefaultQuery("history", "false")
	if history == "true" {
		count++
	}
	allTags := c.DefaultQuery("allTags", "false")
	if allTags == "true" {
		count++
	}

	if count > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "only allow one of these query param at a time (tag / user / follow / self / history / allTags / likedBy)",
		})
		return
	}

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

	if count == 0 {
		/* 公共内容 */
		contents = model.QueryContents("public", "_", orderBy, order, num)
	} else if tag != "" {
		/* 指定tag */
		contents = model.QueryContents("tag", tag, orderBy, order, num)
	} else if username != "" {
		/* 指定user */
		userID, _ := model.QueryUserIDWithName(username)
		contents = model.QueryContents("user", userID, orderBy, order, num)
	} else if follow == "true" {
		/* 我关注的 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("follow", loginUserID, orderBy, order, num)
	} else if self == "true" {
		/* 我的 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("user", loginUserID, orderBy, order, num)
	} else if history == "true" {
		/* 我的浏览记录 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("history", loginUserID, "_", "_", num)
	} else if search != "" {
		/* 搜索 */
		contents = model.QueryContents("search", search, orderBy, order, num)
	} else if allTags == "true" {
		/* 我关注的 tag 的全部内容 */
		// 获得已登录用户的 userID
		loginUserID, err := GetUserIDByAuth(c)
		if err != nil {
			return
		}
		contents = model.QueryContents("allTags", loginUserID, orderBy, order, num)
	} else if likedBy != "" {
		/* 某人喜欢的所有内容 */
		contents = model.QueryContents("like", likedBy, orderBy, order, num)
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
		fmt.Println("contentID (integer) required")
		return
	}

	if err := model.DeleteContentWithContentID(loginUserID, contentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// PostContent 发布一条内容，请求体应该为 Form-data 形式，要求内容有:
// 1. title : 视频标题
// 2. description : 视频详细介绍
// 3. video : 视频文件，支持格式有: WMV,AVI,MKV,RMVB,MP4,MOV; 大小限制: 200 mb
// 4. cover : 封面图片文件，支持格式有: jpg,png,jpeg,svg; 大小限制: 1 mb
// 5. duration : 视频长度，以秒为单位
func PostContent(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	videoFile, err1 := c.FormFile("video")
	coverFile, err2 := c.FormFile("cover")
	durationStr := c.PostForm("duration")
	duration, err3 := strconv.Atoi(durationStr)

	if title == "" || description == "" || err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect Form-data: {title(Text), description(Text), video(File), cover(File), dutation(int)}",
		})
		fmt.Println("expect Form-data: {title(Text), description(Text), video(File), cover(File), dutation(int)}")
		return
	}

	// 检查视频类型
	allowedVideoTypes := []string{".wmv", ".avi", ".mkv", ".rmvb", ".mp4", ".MOV"}
	videoSuffix := path.Ext(videoFile.Filename)

	if !contains(allowedVideoTypes, videoSuffix) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "supported video type: WMV,AVI,MKV,RMVB,MP4 ",
		})
		fmt.Println("supported video type: WMV,AVI,MKV,RMVB,MP4")
		return
	}

	// 检查视频大小
	if videoFile.Size > (200 << 20) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "video file too big (> 200mb)",
		})
		fmt.Println("video file too big (> 200mb)")
		return
	}

	// 检查封面图片类型
	allowedImageTypes := []string{".jpg", ".png", ".jpeg", ".svg"}
	coverSuffix := path.Ext(coverFile.Filename)
	if !contains(allowedImageTypes, coverSuffix) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "supported image type: jpg,png,jpeg,svg ",
		})
		fmt.Println("supported image type: jpg,png,jpeg,svg ")
		return
	}

	// 检查封面图片大小
	if coverFile.Size > (10 << 20) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "image file too big (> 10mb)",
		})

		fmt.Println("image file too big (> 10mb)")
		return
	}

	// 生成文件路径 和 URL
	randomNumberSuffix := rand.Intn(1000)
	predictedContentID := model.QueryMaxContentID() + 1
	videoPath := fmt.Sprintf("/home/lighthouse/IOS_Files/videos/content%d_video%d%s", predictedContentID, randomNumberSuffix, videoSuffix)
	coverPath := fmt.Sprintf("/home/lighthouse/IOS_Files/covers/content%d_cover%d%s", predictedContentID, randomNumberSuffix, coverSuffix)
	videoURL := fmt.Sprintf("/static/videos/content%d_video%d%s", predictedContentID, randomNumberSuffix, videoSuffix)
	coverURL := fmt.Sprintf("/static/covers/content%d_cover%d%s", predictedContentID, randomNumberSuffix, coverSuffix)

	// 保存到服务器
	if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "Save video failed",
		})
		return
	}

	if err := c.SaveUploadedFile(coverFile, coverPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "Save cover image failed",
		})
		return
	}

	// 更新数据库
	if err := model.InsertContent(title, description, coverURL, videoURL, loginUserID, duration); err != nil {
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

type contentTagFormat struct {
	ContentID int    `json:"contentID" binding:"required"`
	Tag       string `json:"tag" binding:"required"`
}

// AddTagForContent : 为内容增加一个 Tag, 要求内容是自己发的
func AddTagForContent(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var newTag contentTagFormat
	if err := c.BindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {contentID, tag}",
		})
		return
	}

	userID, err := model.QueryUserIDWithContentID(newTag.ContentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if userID != loginUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"error":  "no access",
		})
		return
	}

	if err := model.InsertContentTag(newTag.ContentID, newTag.Tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	tags, _ := model.QueryTagsWithContentID(newTag.ContentID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tags,
	})
}

// DeleteTagForContent : 为内容删除一个 Tag, 要求内容是自己发的
func DeleteTagForContent(c *gin.Context) {
	// 获得已登录用户的 userID
	loginUserID, err := GetUserIDByAuth(c)
	if err != nil {
		return
	}

	var newTag contentTagFormat
	if err := c.BindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "expect JSON: {contentID, tag}",
		})
		return
	}

	userID, err := model.QueryUserIDWithContentID(newTag.ContentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if userID != loginUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"error":  "no access",
		})
		return
	}

	if err := model.DeleteContentTag(newTag.ContentID, newTag.Tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	tags, _ := model.QueryTagsWithContentID(newTag.ContentID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tags,
	})
}

// utility func
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
