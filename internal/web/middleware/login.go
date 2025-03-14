package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type MiddlewareBuilder struct{}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
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
				// 不影响系统正常使用的记录日志即可
				fmt.Println(err)
			}
		}
	}
}
