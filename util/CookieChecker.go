package util

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CookieChecker(ctx *gin.Context, key string) bool { //쿠키가 세션에 저장 되어있는지 확인되면 true 없으면 false
	session := sessions.Default(ctx)

	if session.Get(key) == nil {
		return false
	}

	return true
}
