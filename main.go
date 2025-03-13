package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/middleware"
	"github.com/l0slakers/webook/internal/repository"
	"github.com/l0slakers/webook/internal/repository/dao"
	"github.com/l0slakers/webook/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/l0slakers/webook/internal/web"
)

func main() {
	server := gin.Default()

	db, err := gorm.Open(mysql.Open("root@root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic("failed to connect database")
	}

	ud := dao.NewUserDAO(db)
	ur := repository.NewUserService(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
	middleware.RegisterMiddleware(server)

	server.Run(":8080")
}
