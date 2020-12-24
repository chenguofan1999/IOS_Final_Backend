package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 为视频建立文件服务
	router.StaticFS("/static", http.Dir("/home/lighthouse/videos"))

	return router
}
