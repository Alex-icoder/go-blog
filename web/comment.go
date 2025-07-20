package web

import (
	"github.com/gin-gonic/gin"
	"go-blog/base/common"
	"go-blog/base/constant"
	"go-blog/base/httpUtil"
	"go-blog/dao"
	"go-blog/middleware"
	"go-blog/service"
	"strconv"
)

func CreateComment(router *gin.Engine) {
	r := router.Group("/v1/comment")
	r.POST("/create", middleware.JWTAuthMiddleware(), createComment)
}

func FindComment(router *gin.Engine) {
	r := router.Group("/v1/comment")
	//  /find?postId=12
	r.GET("/find", middleware.JWTAuthMiddleware(), findComment)
}

func findComment(c *gin.Context) {
	postIdStr := c.Query(constant.PostId)
	var postId uint
	if len(postIdStr) > 0 {
		post, err := strconv.ParseUint(postIdStr, 10, 64)
		if err != nil {
			common.FailBiz(c, "文章ID参数格式错误")
			return
		}
		postId = uint(post)
	}
	if postId == 0 {
		common.FailBiz(c, "未指定文章")
		return
	}
	list, err := service.FindComment(postId)
	if err != nil {
		common.FailBiz(c, err.Error())
		return
	}
	common.Success(c, list)
}

func createComment(c *gin.Context) {
	userId := c.GetUint(constant.UserId)
	if userId == 0 {
		common.FailBiz(c, "获取用户登录信息失败")
		return
	}
	vo, err := httpUtil.ProcessPostForm[CommentVo](c)
	if err != nil {
		//入参校验合法性
		httpUtil.HandleValidationError(c, err)
		return
	}

	comment := dao.Comment{
		PostID:  vo.PostID,
		Content: vo.Content,
		UserID:  userId,
	}

	errs := service.CreateComment(&comment)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	common.SuccessWithMsg(c, gin.H{constant.CommentId: comment.Model.ID}, "评论创建成功")
}

type CommentVo struct {
	PostID  uint   `form:"postId" binding:"required"`
	Content string `form:"content"  binding:"required,min=2,max=3000"`
}
