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
		var User models.User
		var DbUser models.User

		if err := c.ShouldBindJSON(&User); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패,  요청 오류",
			})

			log.Println(err.Error())
			return

		}

		if err := db.Where("username = ?", User.Username).First(&DbUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "아이디가 중복됩니다."})
			log.Println("Username already exists:", err.Error())
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "서버 오류"})
			log.Println("Database error:", err.Error())
			return
		}

		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패, 암호화 오류",
			})
			log.Println(err.Error())
			return
		}

		User.Password = string(bcryptPassword)
		db.Create(&User)
		c.JSON(http.StatusOK, gin.H{
			"msg": "회원가입 완료",
		})

		return
	}

}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestUser models.User
		var dbUser models.User

		err := c.ShouldBindJSON(&requestUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "입력값 오류"})
			log.Println(err.Error())
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(requestUser.Password)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "아이디 또는 패스워드가 틀렸습니다."})
			return
		}

		session := sessions.Default(c)
		session.Set("username", dbUser.Username)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"msg": "로그인 완료! 쿠키 발급"})
		return
	}

}
