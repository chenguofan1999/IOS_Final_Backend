package controller

import (
	//"fmt"

	"ios/model"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type loginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// used to create and parse token
var serversecret = []byte("randomkey")

// GetNameByToken 由 jwt token 得到用户名，解析失败则返回 “”
func GetNameByToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return serversecret, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if username, ok := claims["username"].(string); ok {
				return username
			}
		}
	}
	return ""
}

// GetTokenByName 由用户名生成 jwt token, 创建失败则返回 “”
func GetTokenByName(name string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": name,
	})
	tokenString, err := token.SignedString(serversecret)
	if err != nil {
		return ""
	}
	return tokenString
}

// SignUp : 注册, 应传入一个对应 loginInfo 的 JSON
func SignUp(c *gin.Context) {
	var info loginInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "binding error",
		})
		return
	}

	if err := model.InsertUser(info.Username, info.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Login : 登录, 应传入一个对应 loginInfo 的 JSON
func Login(c *gin.Context) {

	// 获取请求体
	var info loginInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "binding error",
		})
		return
	}

	// 取得用户密码
	expected, err := model.QueryPasswordWithName(info.Username)
	if err != nil {
		// 用户不存在
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	// 验证密码
	if expected != "" && info.Password == expected {
		tokenString := GetTokenByName(info.Username)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"token":  tokenString,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"error":  "password error",
		})
	}
}
