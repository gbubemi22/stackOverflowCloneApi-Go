package main

import (
	"github.com/gbubemi22/go-stackOverflow/database"
	routes "github.com/gbubemi22/go-stackOverflow/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	//middleware "github.com/gbubemi22/go-stackOverflow/middleware"
	"log"
	"os"
)

var qestionCollection *mongo.Collection = database.OpenCollection(database.Client, "question")

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.QestionRoutes(router)
	routes.AnswerRoutes(router)
	//router.Use(middleware.Authenticate())

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	//Serve the Swagger UI at the /swagger URL
	router.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + port)

}
