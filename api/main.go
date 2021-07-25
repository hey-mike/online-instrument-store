package main

import (
	"context"
	"time"

	"os"

	log "github.com/sirupsen/logrus"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "microservice/docs"
	"microservice/src/config"
	"microservice/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesController *controllers.RecipesController

func init() {
	ctx := context.Background()
	mongo_uri := config.GetEnv("MONGO_URI")

	log.SetLevel(log.DebugLevel)
	log.Debug("MONGO_URI", mongo_uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to MongoDB")

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Info(status)

	recipesController = controllers.NewRecipesController(ctx, collection, redisClient)
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

// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()

	router.POST("/recipes", recipesController.NewRecipeController)
	router.GET("/recipes", recipesController.ListRecipesController)
	router.PUT("/recipes/:id", recipesController.UpdateRecipeController)
	router.DELETE("/recipes/:id", recipesController.DeleteRecipeController)
	router.GET("/recipes/:id", recipesController.GetRecipeController)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
