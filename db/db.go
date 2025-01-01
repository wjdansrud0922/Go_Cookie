package db

import (
	"golangCRUD/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initdb() *gorm.DB {
	dsn := "root:1234@tcp(127.0.0.1:3306)/GO_CRUD?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}
