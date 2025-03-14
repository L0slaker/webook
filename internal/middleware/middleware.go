package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/web/middleware"
	"strings"
	"time"
)

func RegisterMiddleware(server *gin.Engine) []gin.HandlerFunc {
	login := middleware.NewMiddlewareBuilder()
	// Session初始化
	//store := cookie.NewStore([]byte("secret"))
	//memstore.NewStore([]byte("secret"))
	//memcached.NewStore()
	// size 最大空闲连接数
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("U2FsdGVkX1+9aG5tSL0nyB2byBGKpuK0"), // authentication 身份认证
		[]byte("U2FsdGVkX19IDF17ov2HRI/9TlXkROBL")) // encryption 数据加密
	if err != nil {
		panic(err)
	}

	return []gin.HandlerFunc{
		corsMiddleware(),
		sessions.Sessions("ssid", store), login.CheckLogin(),
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
		//ExposeHeaders:              nil,
		MaxAge: 12 * time.Hour,
	})
}
