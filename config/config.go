package config

import (
	"go-blog/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Dsn string
var JwtSecretKey string
var Port string
var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// 自动迁移数据库表
	errs := DB.AutoMigrate(&dao.User{}, &dao.Post{}, &dao.Comment{})
	if errs != nil {
		log.Fatal("自动迁移数据库表失败：", err)
		return
	}
	log.Println("自动迁移数据库表完成")
}
