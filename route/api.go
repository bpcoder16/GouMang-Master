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
	}

	return apiRouter
}
