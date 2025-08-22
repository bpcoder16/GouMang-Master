package route

import (
	"goumang-master/controller"

	"github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"
)

func Api() gin.Router {
	apiRouter := gin.NewDefaultRouter("/api")

	taskGroup := apiRouter.Group("/task")
	{
		taskGroup.GET("/config", (&controller.Task{}).Config)
	}

	return apiRouter
}
