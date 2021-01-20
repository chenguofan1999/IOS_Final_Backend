package routes

import (
	"net/http"

	"ios/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 建立静态文件服务
	router.StaticFS("/static", http.Dir("/home/lighthouse/IOS_Files"))
	router.MaxMultipartMemory = 200 << 20 // 仅用于限制视频大小，图片(头像、视频封面)大小限制为 1MB

	/************ 用户服务 **************/
	// 注册与登录
	router.POST("/signup", controller.SignUp)
	router.POST("/login", controller.Login)

	// 用户信息
	router.GET("/users/:username", controller.GetUserInfoByName)    // 获取某用户详细信息
	router.GET("/user", controller.GetSelfInfo)                     // 获取自己的用户信息
	router.PUT("/user/bio", controller.UpdateUserBio)               // 更新自己的简介
	router.PUT("/user/avatar", controller.UpdateUserAvatar)         // 更新自己的头像
	router.POST("/user/avatar", controller.UpdateUserAvatar)        // 更新自己的头像
	router.GET("/user/tags", controller.GetTagsForCurrentUser)      // 为自己增加关注的 tag
	router.POST("/user/tags", controller.AddTagForCurrentUser)      // 为自己增加关注的 tag
	router.DELETE("/user/tags", controller.DeleteTagForCurrentUser) // 为自己删除关注的 tag

	// 关注
	router.GET("/users/:username/followers", controller.GetFollowersByUserID) // 获取某用户关注者
	router.GET("/users/:username/following", controller.GetFollowingByUserID) // 获取某用户关注的人
	router.GET("/user/following/:username", controller.CheckFollowing)        // 检查是否已关注某用户
	router.PUT("/user/following/:username", controller.FollowUser)            // 关注某用户
	router.DELETE("/user/following/:username", controller.UnfollowUser)       // 取消关注某用户

	/************ Content 服务 **************/
	router.GET("/contents", controller.GetContents)                      // 获取内容集，详见 controller.GetContents 注释
	router.POST("/contents", controller.PostContent)                     // 发布内容
	router.GET("/contents/:contentID", controller.GetContentByContentID) // 根据 contentID 获取某条内容的详细信息
	router.DELETE("/contents/:contentID", controller.DeleteContent)      // 删除内容，仅能删除自己发出的内容
	router.POST("/content/tags", controller.AddTagForContent)            // 为内容增加 tag , 仅能为自己发的内容增加标签
	router.DELETE("/content/tags", controller.DeleteTagForContent)       // 为内容删除 tag , 仅能为自己发的内容删除标签
	// Todo: POST /content

	// Todo: PUT /content/:contentID (maybe)

	/************ Comment 服务 **************/
	router.GET("/comments", controller.GetComments)                 // 获取评论集，详见 controller.GetComments 注释
	router.POST("/comments", controller.PostComment)                // 发布评论
	router.DELETE("/comments/:commentID", controller.DeleteComment) // 删除评论, 仅能删除自己发的评论

	/************ Reply 服务 **************/
	router.GET("/replies", controller.GetReplies)              // 获取回复集，详见 controller.GetReplies 注释
	router.POST("/replies", controller.PostReply)              // 发布回复
	router.DELETE("/replies/:replyID", controller.DeleteReply) // 删除回复, 仅能删除自己发的回复

	/************ Like **************/
	router.PUT("/like/content/:contentID", controller.LikeContent)
	router.DELETE("/like/content/:contentID", controller.CancelLikeContent)
	router.PUT("/like/comment/:commentID", controller.LikeComment)
	router.DELETE("/like/comment/:commentID", controller.CancelLikeComment)
	router.PUT("/like/reply/:replyID", controller.LikeReply)
	router.DELETE("/like/reply/:replyID", controller.CancelLikeReply)

	return router
}
