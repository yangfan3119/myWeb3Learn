package router

import (
	"go01_4/config"
	"go01_4/controllers"
	"go01_4/middlewares"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine) {
	controllers.Init_db()

	userCtr := controllers.NewUserCtr(&config.Cfg)
	postCtr := controllers.NewPostCtr()
	commCtr := controllers.NewCommentCtr()
	// 无需认证的路由
	router.POST("/register", userCtr.Register)
	router.POST("/login", userCtr.Login)

	router.GET("/post/", postCtr.GetPostArray)
	router.GET("/post/:id", postCtr.GetPost)
	router.GET("/post/:id/comment", commCtr.GetComms)

	// 需认证的文章路由
	postR := router.Group("/post").Use(middlewares.JWTCheckUserAuth(&config.Cfg))
	{
		postR.POST("/create", postCtr.Create)
		postR.PUT("/:id", postCtr.Update)
		postR.DELETE("/:id", postCtr.Delpost)
	}
	// 需认证的评论路由
	commentR := router.Group("/comm").Use(middlewares.JWTCheckUserAuth(&config.Cfg))
	{
		commentR.POST("/create", commCtr.Create)
		commentR.PUT("/:id", commCtr.Update)
		commentR.DELETE("/:id", commCtr.DelComm)
	}
}
