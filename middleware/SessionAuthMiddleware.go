package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if err := session.Get("username"); err != true {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "로그인이 필요합니다."})
			return
		}

		c.Next()
	}
}
