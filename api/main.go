package main

import (
	"context"
	"log"
	"os"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "microservice/docs"
	"microservice/src/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesController *controllers.RecipesController

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesController = controllers.NewRecipesController(ctx, collection)
}

// @title Recipe API
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

	router.POST("/recipes", recipesController.NewRecipeController)
	router.GET("/recipes", recipesController.ListRecipesController)
	router.PUT("/recipes/:id", recipesController.UpdateRecipeController)
	router.DELETE("/recipes/:id", recipesController.DeleteRecipeController)
	router.GET("/recipes/:id", recipesController.GetRecipeController)

	// url := ginSwagger.URL("http://petstore.swagger.io:8080/swagger/doc.json") // The url pointing to API definition
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
