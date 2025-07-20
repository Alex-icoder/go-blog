package web

import (
	"github.com/gin-gonic/gin"
	"go-blog/base/common"
	"go-blog/base/constant"
	"go-blog/base/httpUtil"
	"go-blog/dao"
	"go-blog/middleware"
	"go-blog/service"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreatePost(router *gin.Engine) {
	r := router.Group("/v1/post")
	r.POST("/create", middleware.JWTAuthMiddleware(), create)
}

func FindPost(router *gin.Engine) {
	r := router.Group("/v1/post")
	r.GET("/find/:id", find)
}

func UpdatePost(router *gin.Engine) {
	r := router.Group("/v1/post")
	r.PUT("/update/:id", middleware.JWTAuthMiddleware(), update)
}

func DeletePost(router *gin.Engine) {
	r := router.Group("/v1/post")
	r.DELETE("/delete/:id", middleware.JWTAuthMiddleware(), delete)
}

func delete(c *gin.Context) {
	userId := c.GetUint(constant.UserId)
	if userId == 0 {
		common.FailBiz(c, "获取用户登录信息失败")
		return
	}
	idStr := c.Param("id")
	var postId uint
	if len(idStr) > 0 {
		post, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			common.FailBiz(c, "文章ID参数格式错误")
			return
		}
		postId = uint(post)
	}
	list, err := service.FindPost(postId)
	if err != nil {
		common.FailBiz(c, err.Error())
		return
	}
	if len(list) == 0 {
		common.FailBiz(c, "文章不存在")
		return
	}
	if list[0].UserID != userId {
		common.FailBiz(c, "只有文章作者自己可以删除")
		return
	}
	errs := service.DeletePost(postId)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	common.SuccessWithMsg(c, nil, "文章删除成功")
}

func update(c *gin.Context) {
	userId := c.GetUint(constant.UserId)
	if userId == 0 {
		common.FailBiz(c, "获取用户登录信息失败")
		return
	}

	idStr := c.Param("id")
	var postId uint
	if len(idStr) > 0 {
		post, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			common.FailBiz(c, "文章ID参数格式错误")
			return
		}
		postId = uint(post)
	}
	list, err := service.FindPost(postId)
	if err != nil {
		common.FailBiz(c, err.Error())
		return
	}
	if len(list) == 0 {
		common.FailBiz(c, "文章不存在")
		return
	}
	if list[0].UserID != userId {
		common.FailBiz(c, "只有文章作者自己可以修改")
		return
	}
	vo, err := httpUtil.ProcessPostJson[PostVo](c)
	if err != nil {
		//入参校验合法性
		httpUtil.HandleValidationError(c, err)
		return
	}
	update := dao.Post{
		Model:   gorm.Model{ID: postId},
		Title:   vo.Title,
		Content: vo.Content,
		UserID:  userId,
	}
	errs := service.UpdatePost(update)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	common.SuccessWithMsg(c, nil, "文章修改成功")
}

func find(c *gin.Context) {
	idStr := c.Param("id")
	var postId uint
	if len(idStr) > 0 {
		post, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			common.FailBiz(c, "文章ID参数格式错误")
			return
		}
		postId = uint(post)
	}
	list, err := service.FindPost(postId)
	if err != nil {
		common.FailBiz(c, err.Error())
		return
	}
	common.Success(c, list)
}

func create(c *gin.Context) {
	userId := c.GetUint(constant.UserId)
	if userId == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{Code: constant.FAIL, Message: "获取用户登录信息失败"})
		return
	}
	vo, err := httpUtil.ProcessPostJson[PostVo](c)
	if err != nil {
		//入参校验合法性
		httpUtil.HandleValidationError(c, err)
		return
	}

	post := dao.Post{
		UserID:  userId,
		Title:   vo.Title,
		Content: vo.Content,
	}
	errs := service.CreatePost(&post)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	common.SuccessWithMsg(c, gin.H{constant.PostId: post.Model.ID}, "文章创建成功")
}

type PostVo struct {
	Title   string `json:"title"  binding:"required,min=2,max=20"`
	Content string `json:"content"  binding:"required,min=2,max=3000"`
}
