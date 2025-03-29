package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/pkg/ginx/rateLimit"
	"github.com/l0slakers/webook/internal/web/middleware"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func RegisterMiddleware(server *gin.Engine) []gin.HandlerFunc {
	login := middleware.NewLoginMiddlewareBuilder()
	// Session初始化
	store := cookie.NewStore([]byte("secret"))
	//memstore.NewStore([]byte("secret"))
	//memcached.NewStore()
	// size 最大空闲连接数
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("U2FsdGVkX1+9aG5tSL0nyB2byBGKpuK0"), // authentication 身份认证
	//	[]byte("U2FsdGVkX19IDF17ov2HRI/9TlXkROBL")) // encryption 数据加密
	//if err != nil {
	//	panic(err)
	//}

	return []gin.HandlerFunc{
		corsMiddleware(),
		sessions.Sessions("ssid", store), login.CheckLoginJWT(),
		rateLimitMiddleware(),
	}
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "my_company.com")
		},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		// 允许前端访问后端响应带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		MaxAge:        12 * time.Hour,
	})
}

func rateLimitMiddleware() gin.HandlerFunc {
	client := redis.NewClient(&redis.Options{
		// TODO 配置文件获取
		Addr: "localhost:6379",
	})

	// TODO 配置文件获取
	return rateLimit.NewBuilder(client, time.Second, 100).Build()
}
