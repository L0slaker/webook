package web

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	user := server.Group("/api/v1/user")

	user.POST("/signup", h.SignUp)
	//user.POST("/login", h.Login)
	user.POST("/login", h.LoginJWT)
	user.POST("/edit", h.Edit)
	user.GET("/info", h.Info)
}
