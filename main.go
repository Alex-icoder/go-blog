package main

import (
	"go-blog/config"
	"go-blog/env"
	"go-blog/router"
	"log"
)

func main() {
	//加载配置
	env.Config()
	//初始化数据库
	config.InitDB()
	//设置路由
	r := router.SetupRouter()
	//启动服务
	err := r.Run(":" + config.Port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
