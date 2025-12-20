package middleware

import (
	"baihu/internal/constant"
	"baihu/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(constant.CookieName)
		if err != nil || token == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 验证 token
		userID, username, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}

// SetAuthCookie 设置认证 Cookie
func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie(constant.CookieName, token, constant.CookieMaxAge, "/", "", false, true)
}

// ClearAuthCookie 清除认证 Cookie
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie(constant.CookieName, "", -1, "/", "", false, true)
}
