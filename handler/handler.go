package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golangCRUD/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

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
	}

}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestUser models.User
		var DbUser models.User

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
	}

}
