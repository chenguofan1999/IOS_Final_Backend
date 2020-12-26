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

	// 关注
	router.GET("/users/:username/followers", controller.GetFollowersByUserID) // 获取某用户关注者
	router.GET("/users/:username/following", controller.GetFollowingByUserID) // 获取某用户关注的人
	router.GET("/user/following/:username", controller.CheckFollowing)        // 当前用户取消关注某用户
	router.PUT("/user/following/:username", controller.FollowUser)            // 当前用户关注某用户
	router.DELETE("/user/following/:username", controller.UnfollowUser)       // 当前用户取消关注某用户

	/************ Content 服务 **************/
	router.GET("/contents", controller.GetContents)                     // 获取内容集，详见 controller.GetContents 注释
	router.GET("/content/:contentID", controller.GetContentByContentID) // 根据 contentID 获取某条内容的详细信息
	return router
}
