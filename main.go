package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/l0slakers/webook/internal/middleware"
	"github.com/l0slakers/webook/internal/repository"
	"github.com/l0slakers/webook/internal/repository/cache"
	"github.com/l0slakers/webook/internal/repository/dao"
	"github.com/l0slakers/webook/internal/service"
	"github.com/l0slakers/webook/internal/web"
)

func main() {
	server := gin.Default()

	db := initDB()
	client := initCache()

	server.Use(middleware.RegisterMiddleware(server)...)

	initUserHandler(db, client, server)

	server.Run(":8080")
}

func initUserHandler(db *gorm.DB, client redis.Cmdable, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(client)
	ur := repository.NewUserRepository(ud, uc)
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

func initCache() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
