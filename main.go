package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/l0slakers/webook/internal/repository"
	"github.com/l0slakers/webook/internal/repository/dao"
	"github.com/l0slakers/webook/internal/service"
	"github.com/l0slakers/webook/internal/web"
)

func main() {
	server := gin.Default()

	db := initDB()

	server.Use(middleware.RegisterMiddleware(server)...)

	initUserHandler(db, server)

	server.Run(":8080")
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserService(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic("failed to connect database")
	}
	err = dao.InitTables(db)
	if err != nil {
		panic("failed to init tables")
	}
	return db
}
