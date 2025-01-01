package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"-"`
}

func main() {
	dsn := "root:1234@tcp(127.0.0.1:3306)/GO_CRUD?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{})

	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.POST("/register", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패 요청 오류",
			})

			log.Println(err.Error())
			return

		}

		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패, 암호화 오류",
			})

			return
		}

		user.Password = string(bcryptPassword)
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"msg": "회원가입 완료",
		})

		return
	})

	router.POST("/login", func(c *gin.Context) {
		var requestUser User
		var DbUser User

		err := c.ShouldBindJSON(&requestUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "입력값 오류"})
			return
		}

		if err := db.Where("username = ?", requestUser.Username).First(&DbUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "해당 유저이름이 존재하지 않습니다."})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(requestUser.Password)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "패스워드가 틀렸습니다."})
			return
		}

		session := sessions.Default(c)
		session.Set("username", DbUser.Username)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"msg": "로그인 완료! 쿠키 발급"})
		return
	})
	router.Run(":8080")
}
