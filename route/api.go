package route

import "github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"

func Api() gin.Router {
	apiRouter := gin.NewDefaultRouter("/api")

	return apiRouter
}
