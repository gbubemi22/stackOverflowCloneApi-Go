package routes

import (
	controller "github.com/gbubemi22/go-stackOverflow/controllers"
	"github.com/gin-gonic/gin"
)

func QestionRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/questions", controller.CreatQestion())
	incomingRoutes.GET("/questions/:question_id", controller.GetOneQestion())
	incomingRoutes.GET("/questions", controller.GetAllQuestion())
	incomingRoutes.PATCH("/questions/:question_id", controller.UpdateQestion())
	incomingRoutes.PUT("/questions/:user_id/:question_id/likes", controller.UpdateLikes())
}

// /:user_id/:question_id/likes
