package routes

import (
	controller "github.com/gbubemi22/go-stackOverflow/controllers"
	"github.com/gin-gonic/gin"
)

func AnswerRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/answers", controller.CreatAnswer())
	incomingRoutes.GET("/answers", controller.GetOneAnswer())
	incomingRoutes.GET("/answers/answer_id", controller.GetOneAnswer())
	incomingRoutes.PATCH("/answers/answer_id", controller.UpdateAnswer())
}
