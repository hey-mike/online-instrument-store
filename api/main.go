package main

import (
	"net/http"
	"time"

	"microservice/src/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

var recipes []models.Recipe

func NewReciptHandler(c *gin.Context) {
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": err.Error()})

		return

	}

	recipe.ID = xid.New().String()

	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.POST("/recipes", NewReciptHandler)

	http.ListenAndServe(":8080", router)
}
