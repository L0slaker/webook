package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MiddlewareBuilder struct{}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 登录校验需要去除“注册”和“登录”两个接口
		if ctx.Request.URL.Path == "/api/v1/user/signup" || ctx.Request.URL.Path == "/api/v1/user/login" {
			return
		}
		sess := sessions.Default(ctx)
		// 用户未登录
		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
