package controllers

import (
	"context"
	"encoding/json"
	"log"
	"microservice/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesController struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesController(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesController {
	return &RecipesController{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// @Summary Returns list of recipes
// @Description get recipes
// @ID get-recipes
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Recipe
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /recipes [get]
func (controller *RecipesController) ListRecipesController(c *gin.Context) {
	val, err := controller.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Printf("Load data from MongoDB")
		cur, err := controller.collection.Find(c, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(c)

		recipes := make([]models.Recipe, 0)
		for cur.Next(controller.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}

		data, _ := json.Marshal(recipes)
		controller.redisClient.Set("recipes", string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Load data from Redis - cache")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}

}

// @Summary Create a new recipe
// @Description create a new recipe
// @ID create-recipe
// @Accept  json
// @Produce  json
// @Param message body models.Recipe true "Recipe Info"
// @Success 200 {object} models.Recipe
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /recipes [post]
func (controller *RecipesController) NewRecipeController(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := controller.collection.InsertOne(controller.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	log.Println("Remove data from Redis")
	controller.redisClient.Del("recipes")

	c.JSON(http.StatusOK, recipe)
}

// @Summary Update a recipe
// @Description update recipe
// @ID update-recipe
// @Accept  json
// @Produce  json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /recipes [put]
func (controller *RecipesController) UpdateRecipeController(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := controller.collection.UpdateOne(controller.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// @Summary Delete a recipe
// @Description delete recipe by ID
// @ID get-recipe
// @Accept  json
// @Produce  json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /recipes/{id} [delete]
func (controller *RecipesController) DeleteRecipeController(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := controller.collection.DeleteOne(controller.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// @Summary Get a recipe
// @Description get string by ID
// @ID get-recipe
// @Accept  json
// @Produce  json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /recipes/{id} [get]
func (controller *RecipesController) GetRecipeController(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := controller.collection.FindOne(controller.ctx, bson.M{
		"_id": objectId,
	})
	var recipe models.Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}
