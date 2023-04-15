package routes

import (
	controllers "github.com/gbubemi22/go-stackOverflow/controllers"
	"github.com/gin-gonic/gin"
	//"github.com/gbubemi22/go-stackOverflow/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/user_id", controller.GetUser())
}