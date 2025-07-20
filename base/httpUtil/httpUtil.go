package httpUtil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-blog/base/common"
	"go-blog/base/constant"
	"net/http"
)

func HandleValidationError(c *gin.Context, err error) {
	// 若不是校验错误，直接返回原始错误
	if _, ok := err.(*validator.InvalidValidationError); ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code:    constant.FAIL,
			Message: "参数验证失败",
		})
		return
	}

	// 处理校验错误
	details := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		tag := err.Tag()
		message := getErrorMessage(fieldName, tag)
		details[fieldName] = message
	}

	c.JSON(http.StatusBadRequest, common.Response{
		Code:    constant.FAIL,
		Message: "参数格式不正确",
		Data:    details,
	})
}

// 根据字段和校验标签获取友好错误信息
func getErrorMessage(field string, tag string) string {
	// 这里可以添加更多的映射关系
	messages := map[string]map[string]string{
		"Email": {
			"required": "邮箱不能为空",
			"email":    "邮箱格式不正确",
			"min":      "邮箱长度不能少于5个字符",
			"max":      "邮箱长度不能超过50个字符",
		},
		"Password": {
			"required": "密码不能为空",
			"min":      "密码长度不能少于6个字符",
			"max":      "",
		},
		"Username": {
			"required": "用户名不能为空",
			"min":      "用户名长度不能少于2个字符",
			"max":      "用户名长度不能超过20个字符",
		},
		"title": {
			"required": "文章标题不能为空",
			"min":      "文章标题长度不能少于2个字符",
			"max":      "文章标题长度不能超过20个字符",
		},
		"content": {
			"required": "文章内容不能为空",
			"min":      "文章内容长度不能少于2个字符",
			"max":      "文章内容长度不能超过3000个字符",
		},
	}

	// 通用错误信息
	defaultMessages := map[string]string{
		"required": "该字段不能为空",
		"min":      "长度不足",
		"max":      "长度超出限制",
		"email":    "邮箱格式不正确",
		"url":      "URL格式不正确",
	}

	if fieldMessages, exists := messages[field]; exists {
		if msg, ok := fieldMessages[tag]; ok {
			return msg
		}
	}

	if msg, exists := defaultMessages[tag]; exists {
		return msg
	}

	return fmt.Sprintf("%s字段格式不正确", field)
}

// 处理表单请求的通用方法
func HandleFormRequest[T any](c *gin.Context) (T, error) {
	var form T
	if err := c.ShouldBind(&form); err != nil {
		return form, err
	}
	return form, nil
}

// 处理POST表单
func ProcessPostForm[T any](c *gin.Context) (T, error) {
	vo, err := HandleFormRequest[T](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return vo, err
	}
	return vo, nil
}

func ProcessPostJson[T any](c *gin.Context) (T, error) {
	var form T
	if err := c.ShouldBindJSON(&form); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return form, err
	}
	return form, nil
}
