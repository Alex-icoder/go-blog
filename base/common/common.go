package common

import (
	"github.com/gin-gonic/gin"
	"go-blog/base/constant"
	"net/http"
)

// 标准化响应格式
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: constant.Success, Data: data})
}

func SuccessWithMsg(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{Code: constant.Success, Data: data, Message: message})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message})
}

func FailBiz(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{Code: constant.FAIL, Message: message})
}
