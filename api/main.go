package main

import (
	"net/http"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "microservice/docs"
	"microservice/src/controllers"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	router := gin.Default()
	// c := controllers.NewRecipeController()

	// v1 := router.Group("/api/v1")
	// {
	// 	recipes := v1.Group("/recipes")
	// 	{
	// 		router.POST("/recipes", c.NewRecipeController)
	// 		router.GET("/recipes", c.ListRecipesController)
	// 		router.PUT("/recipes/:id", c.UpdateRecipeController)
	// 		router.DELETE("/recipes/:id", c.DeleteRecipeController)
	// 		router.GET("/recipes/:id", c.GetRecipeController)
	// 	}
	// 	//...
	// }
	// router.GET("/swagger/*any", ginSwagger.WrapController(swaggerFiles.Controller))

	router.POST("/recipes", controllers.NewRecipeController)
	router.GET("/recipes", controllers.ListRecipesController)
	router.PUT("/recipes/:id", controllers.UpdateRecipeController)
	router.DELETE("/recipes/:id", controllers.DeleteRecipeController)
	router.GET("/recipes/:id", controllers.GetRecipeController)

	// url := ginSwagger.URL("http://petstore.swagger.io:8080/swagger/doc.json") // The url pointing to API definition
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
