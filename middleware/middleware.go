package middleware

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-blog/base/common"
	"go-blog/base/constant"
	"go-blog/config"
	"log"
	"net/http"
	"strings"
	"time"
)

// 耗时统计中间件
func LatencyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set(constant.RequestID, uuid.NewString())
		c.Next()
		latency := time.Since(start)
		log.Printf("requestID:%s,method:%s,path:%s,cost:%v", c.GetString(constant.RequestID), c.Request.Method, c.Request.URL.Path, latency)
	}
}

// 跨域中间件配置
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// RBAC权限中间键
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		for _, r := range roles.([]string) {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

// 鉴权中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{Code: constant.FAIL, Message: "用户未登录"})
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JwtSecretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{Code: constant.FAIL, Message: err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 注意：JWT 中的数字会被解析为 float64，需要转换回 uint
			if userIdFloat, ok := claims[constant.UserId].(float64); ok {
				userId := uint(userIdFloat)
				c.Set(constant.UserId, userId)
			}
			c.Set(constant.UserName, claims[constant.UserName])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{Code: constant.FAIL, Message: "无效的token"})
		}
	}
}
