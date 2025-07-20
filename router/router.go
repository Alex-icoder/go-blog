package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/middleware"
	"go-blog/web"
)

func SetupRouter() *gin.Engine {
	// 新建一个没有任何默认中间件的路由
	router := gin.New()
	router.Use(
		// 全局中间件
		// Logger 中间件将日志写入 gin.DefaultWriter，即使你将 GIN_MODE 设置为 release。
		// By default gin.DefaultWriter = os.Stdout
		gin.Logger(),
		middleware.LatencyLogger(),
		// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
		gin.Recovery(),
		middleware.CORSMiddleware(),
	)
	web.UserRegister(router)
	web.UserLogin(router)

	web.CreatePost(router)
	web.FindPost(router)
	web.UpdatePost(router)
	web.DeletePost(router)

	web.CreateComment(router)
	web.FindComment(router)
	return router
}
