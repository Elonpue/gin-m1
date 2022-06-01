package router

import (
	"gin-m1/controller"
	"gin-m1/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	router := gin.New()
	//宕机时可以恢复
	router.Use(gin.Recovery())

	//跨域中间件
	router.Use(middleware.Cors())
	//日志中间件
	router.Use(middleware.Logger())

	register(router)

	return router
}

func register(r *gin.Engine) {
	userController := controller.NewUserController()
	r.POST("/api/login", userController.Login)
	r.GET("/api/menus", middleware.AuthMiddleware(), controller.GetMenus)
	r.GET("/api/users", middleware.AuthMiddleware(), userController.GetUsers)
	r.POST("/api/users", middleware.AuthMiddleware(), userController.CreateUser)
	r.PUT("/api/users/:id", userController.UpdateUser)
	r.DELETE("/api/users/:id", middleware.AuthMiddleware(), userController.DeleteUser)
}
