package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-blog/base/common"
	"go-blog/base/constant"
	"go-blog/base/httpUtil"
	"go-blog/config"
	"go-blog/dao"
	"go-blog/service"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func UserRegister(router *gin.Engine) {
	r := router.Group("/v1/user")
	r.POST("/register", register)
}

func UserLogin(router *gin.Engine) {
	r := router.Group("/v1/user")
	r.POST("/login", login)
}

func login(c *gin.Context) {
	var vo LoginVO
	vo, err := httpUtil.ProcessPostJson[LoginVO](c)
	if err != nil {
		//入参校验合法性
		httpUtil.HandleValidationError(c, err)
		return
	}

	user, errs := service.UserLogin(vo.Username, vo.Password)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	//生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constant.UserId:   user.ID,
		constant.UserName: user.Username,
		constant.Exp:      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		log.Println(err)
		common.FailBiz(c, "登录失败")
		return
	}
	common.SuccessWithMsg(c, gin.H{"token": tokenString}, "登录成功")
}

func register(c *gin.Context) {
	vo, err := httpUtil.ProcessPostJson[RegisterVO](c)
	if err != nil {
		//入参校验合法性
		httpUtil.HandleValidationError(c, err)
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(vo.Password), bcrypt.DefaultCost)
	if err != nil {
		common.FailBiz(c, "注册失败:密码加密失败")
		return
	}
	user := dao.User{
		Username: vo.Username,
		Password: string(hashedPassword),
		Email:    vo.Email,
	}
	errs := service.UserRegister(user)
	if errs != nil {
		common.FailBiz(c, errs.Error())
		return
	}
	common.Success(c, "注册成功")
}

type RegisterVO struct {
	Username string `json:"username" form:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
	Email    string `json:"email"  form:"email"    binding:"required,email,min=5,max=50"`
}

type LoginVO struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}
