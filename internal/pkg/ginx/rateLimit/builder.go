package rateLimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

//go:embed slide_window.lua
var luaScript string

type Builder struct {
	prefix   string
	cmd      redis.Cmdable
	interval time.Duration // 窗口大小
	rate     int           // 允许通过的请求数量
}

func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		prefix:   "ip-rate-limiter",
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isLimit, err := b.limit(ctx)
		if err != nil {
			//TODO 日志
			log.Println(err)
			// 这里可以分为两种做法：激进和保守
			// 激进：不用影响我系统的正常使用，不进行限流
			// 保守：redis崩溃了，这里我们要保护系统的话必须对整个系统限流
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if isLimit {
			//TODO 日志
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
