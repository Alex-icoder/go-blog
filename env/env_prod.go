//go:build prod
// +build prod

package env

import (
	"go-blog/config"
	"log"
)

func Config() {
	log.Println("env:prod")
	config.Dsn = "root:root@admin@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local"
	config.JwtSecretKey = "myblog"
	config.Port = "8090"
}
