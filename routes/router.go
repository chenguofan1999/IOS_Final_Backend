package routes

import (
	"net/http"

	"ios/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 为视频建立文件服务
	router.StaticFS("/static", http.Dir("/home/lighthouse/videos"))

	/************ 用户服务 **************/
	// 注册与登录
	router.POST("/signup", controller.SignUp)
	router.POST("/login", controller.Login)

	// 用户信息
	router.GET("/users/:username", controller.GetUserInfoByName) // 获取某用户详细信息
	router.GET("/user", controller.GetSelfInfo)
	// Todo: 修改 bio
	// Todo: 修改 avatar
	// Todo: 增加 tag

	// 关注
	router.GET("/users/:username/followers", controller.GetFollowersByUserID) // 获取某用户关注者
	router.GET("/users/:username/following", controller.GetFollowingByUserID) // 获取某用户关注的人
	router.GET("/user/following/:username", controller.CheckFollowing)        // 检查是否已关注某用户
	router.PUT("/user/following/:username", controller.FollowUser)            // 关注某用户
	router.DELETE("/user/following/:username", controller.UnfollowUser)       // 取消关注某用户

	/************ Content 服务 **************/
	router.GET("/contents", controller.GetContents)                      // 获取内容集，详见 controller.GetContents 注释
	router.GET("/contents/:contentID", controller.GetContentByContentID) // 根据 contentID 获取某条内容的详细信息
	// Todo: DELETE /content/:contentID
	// Todo: PUT /content/:contentID (maybe)
	// Todo: POST /content

	/************ Comment 服务 **************/
	router.GET("/comments", controller.GetComments)                 // 获取评论集，详见 controller.GetComments 注释
	router.POST("/comments", controller.PostComment)                // 发布评论
	router.DELETE("/comments/:commentID", controller.DeleteComment) // 删除评论, 仅能删除自己发的评论

	/************ Reply 服务 **************/
	router.GET("/replies", controller.GetReplies)              // 获取回复集，详见 controller.GetReplies 注释
	router.POST("/replies", controller.PostReply)              // 发布回复
	router.DELETE("/replies/:replyID", controller.DeleteReply) // 删除回复, 仅能删除自己发的回复

	/************ Like **************/
	// Todo: PUT /like/content/:contentID
	// Todo: DELETE /like/content/:contentID
	// Todo: PUT /like/comment/:commentID
	// Todo: DELETE /like/comment/:commentID
	// Todo: PUT /like/reply/:replyID
	// Todo: DELETE /like/reply/:replyID

	return router
}
