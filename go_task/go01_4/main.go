package main

import (
	"go01_4/config"
	"go01_4/databases"
	"go01_4/logger"
	"go01_4/middlewares"
	"go01_4/models"
	"go01_4/router"

	"github.com/gin-gonic/gin"
)

func testInit() {
	config.Cfg.SqliteName = ":memory:"
	db := databases.GetDB()
	// 添加3组用户、文章、评论测试数据
	users := []models.User{
		{Username: "user1", Password: "$2a$10$SJNHwr5mrmcD5S5g887Z3.OX3fhw4dnd.XhhUEaOR8wW3V1ot8aN", Email: "user1@test.com"},
		{Username: "user2", Password: "$2a$10$/J4gz.aheRZc6vLI2HrgIOPyG7UqmN/ndWUBOJ4HnILOdjWSNuNNC", Email: "user2@test.com"},
		{Username: "user3", Password: "$2a$10$39OzxaUhP5Sw4GBligDpc.ebAhQc77shI4Cw5uIGI9tmIzDVxw.2y", Email: "user3@test.com"},
	}

	for i := range users {
		db.Create(&users[i])
	}
	posts := []models.Post{
		{Title: "Post1", Content: "Content1", UserID: 1},
		{Title: "Post2", Content: "Content2", UserID: 2},
		{Title: "Post3", Content: "Content3", UserID: 3},
	}
	for i := range posts {
		db.Create(&posts[i])
	}
	comments := []models.Comment{
		{Content: "Comment1", UserID: 1, PostID: 1},
		{Content: "Comment2", UserID: 2, PostID: 2},
		{Content: "Comment3", UserID: 3, PostID: 3},
	}
	for i := range comments {
		db.Create(&comments[i])
	}
}

func main() {
	// 配置文件读取与解析
	config.Cfg.Load("./config/dev.config.yaml")
	// 日志初始化
	logger.Init()

	// 数据库初始化
	databases.InitDB(config.Cfg.SqliteName)
	if databases.GetDB() == nil {
		return
	}
	models.MigrateDB(databases.GetDB())
	testInit()

	// 路由初始化
	r := gin.New()
	r.Use(middlewares.RecoveryMiddleware())
	router.SetRouter(r)

	// 服务启动，端口监听
	r.Run(":" + config.Cfg.ServerPort)
}
