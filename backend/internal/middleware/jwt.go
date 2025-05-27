package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte("hi-wine-jwt-secret") // 建议放配置

func JwtAuth(r *ghttp.Request) {
	// 1. 获取 Authorization 头部
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "未登录或token缺失",
		})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. 解析和校验 JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil || !token.Valid {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "token无效或已过期",
		})
		return
	}

	// 3. 将用户信息写入 context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		r.SetCtxVar("userId", claims["userId"])
		r.SetCtxVar("openid", claims["openid"])
	}
	r.Middleware.Next()
}
