package handlers

import (
	"context"
	"net/http"
	"tutorial/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeHandler struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewRecipeHandler(ctx context.Context, collection mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		ctx:        ctx,
		collection: &collection,
	}
}

func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	curr, err := handler.collection.Find(handler.ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	defer curr.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)
	for curr.Next(handler.ctx) {
		var recipe models.Recipe
		curr.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}
