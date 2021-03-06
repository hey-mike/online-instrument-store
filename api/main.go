package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "microservice/docs"
	"microservice/src/controllers"
	"microservice/src/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesController *controllers.RecipesController
var authController *controllers.AuthController

func init() {
	godotenv.Load()
	ctx := context.Background()

	mongo_uri := os.Getenv("MONGO_URI")
	mongo_db := os.Getenv("MONGO_DATABASE")

	log.SetLevel(log.DebugLevel)
	log.Debug("MONGO_URI: ", mongo_uri)
	log.Debug("MONGO_DATABASE: ", mongo_db)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_uri))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to MongoDB")

	collection := client.Database(mongo_db).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Info("Connected to Redis: ", status)

	recipesController = controllers.NewRecipesController(ctx, collection, redisClient)

	collectionUsers := client.Database(mongo_db).Collection("users")
	authController = controllers.NewAuthController(ctx, collectionUsers)

}
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("API_VERSION")})
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

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func SetupServer() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middlewares.RequestIdMiddleware())
	router.Use(middlewares.PrometheusMiddleware())

	router.POST("/signin", authController.SignIn)
	router.GET("/recipes", recipesController.ListRecipes)
	router.POST("/refresh", authController.RefreshToken)

	authorized := router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesController.NewRecipe)
		authorized.PUT("/recipes/:id", recipesController.UpdateRecipe)
		authorized.DELETE("/recipes/:id", recipesController.DeleteRecipe)
		router.GET("/recipes/:id", recipesController.GetRecipe)
	}
	router.GET("/version", VersionHandler)
	router.GET("/prometheus", gin.WrapH(promhttp.Handler()))

	// enable swagger doc
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router;
}

func main() {
	if err := SetupServer().Run(":8000"); err != nil {
        log.Fatal("Failed to Run Server")
    }
}
