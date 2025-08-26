package route

import (
	"goumang-master/controller"

	"github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"
)

func Api() gin.Router {
	apiRouter := gin.NewDefaultRouter("/api")

	taskGroup := apiRouter.Group("/task")
	{
		taskGroup.GET("/list", (&controller.Task{}).List)
		taskGroup.GET("/config", (&controller.Task{}).Config)
		taskGroup.POST("/create", (&controller.Task{}).Create)
		taskGroup.GET("/detail", (&controller.Task{}).Detail)
		taskGroup.POST("/immediately-run", (&controller.Task{}).ImmediatelyRun)
		taskGroup.POST("/delete", (&controller.Task{}).Delete)
		taskGroup.POST("/enable", (&controller.Task{}).Enable)
		taskGroup.POST("/disable", (&controller.Task{}).Disable)
	}

	return apiRouter
}
