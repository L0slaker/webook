package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/l0slakers/webook/internal/web"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct{}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// error gob: type not registered for interface: time.Time
	gob.Register(time.Now()) // redis注册类型

	return func(ctx *gin.Context) {
		// 登录校验需要去除“注册”和“登录”两个接口
		if ctx.Request.URL.Path == "/api/v1/user/signup" || ctx.Request.URL.Path == "/api/v1/user/login" {
			return
		}
		sess := sessions.Default(ctx)
		// 用户未登录
		userId := sess.Get("userId")
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新登录态
		// 默认策略（实际情况可以以产品为准）：这边默认一分钟刷新一次
		// TODO 通过配置文件获取
		now := time.Now()
		const updateTimeKey = "update_time"
		// 获取上一次更新时间比较
		val := sess.Get(updateTimeKey)
		last, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(last) > time.Minute {
			sess.Set(updateTimeKey, now)
			// session自身的问题
			// 如果不重新设置，就会用update_time覆盖掉user_id，届时user_id这个key就不存在了
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// TODO 记录日志
				// 不影响系统正常使用的记录日志即可
				fmt.Println(err)
			}
		}
	}
}

func (m *LoginMiddlewareBuilder) CheckLoginJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 登录校验需要去除“注册”和“登录”两个接口
		if ctx.Request.URL.Path == "/api/v1/user/signup" || ctx.Request.URL.Path == "/api/v1/user/login" {
			return
		}
		// 用户未登录
		tokenStr := ctx.GetHeader("x-jwt-token")
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//segs := strings.Split(" ", authCode)
		//if len(segs) != 2 {
		//	// 传入的 Authentication 不正确
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		//tokenStr := segs[1]
		var claim web.UserClaim
		token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JwtKey), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 校验token是否为空&是否过期
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新登录态：剩余过期时间<50s
		if claim.ExpiresAt.Sub(time.Now()) < time.Second*50 {
			claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 15))
			tokenStr, err = token.SignedString([]byte(web.JwtKey))
			if err != nil {
				// TODO 记录日志
				// 无需中断，仅仅是过期，不影响系统正常使用
				fmt.Println(err)
			}
			// 写入头部返回给前端
			ctx.Header("x-jwt-token", tokenStr)
		}

		// 方便其他功能使用
		ctx.Set("user", claim)
	}
}
