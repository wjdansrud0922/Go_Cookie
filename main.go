package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golangCRUD/db"
	"golangCRUD/handler"
)

// TODO: Add .env
//
//	AppState
func main() {

	db := db.Initdb()
	store := cookie.NewStore([]byte("secret"))

	router := gin.Default()

	router.Use(sessions.Sessions("mysession", store))

	router.POST("/register", handler.RegisterHandler(db))
	router.POST("/login", handler.LoginHandler(db))

	router.Run(":8080")
}
