package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/web/middleware"
	"strings"
	"time"
)

func RegisterMiddleware(server *gin.Engine) []gin.HandlerFunc {
	login := middleware.NewMiddlewareBuilder()
	// Session初始化
	store := cookie.NewStore([]byte("secret"))

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
