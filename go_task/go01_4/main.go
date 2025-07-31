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

	// 路由初始化
	r := gin.New()
	r.Use(middlewares.RecoveryMiddleware())
	router.SetRouter(r)

	// 服务启动，端口监听
	r.Run(":" + config.Cfg.ServerPort)
}
