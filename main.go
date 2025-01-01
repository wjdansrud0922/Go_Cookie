package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golangCRUD/db"
	"golangCRUD/handler"
)

func main() {

	db := db.Initdb()
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.POST("/register", handler.RegisterHandler(db))
	router.POST("/login", handler.LoginHandler(db))
	router.Run(":8080")
}
