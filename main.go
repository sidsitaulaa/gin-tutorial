// Recipes API
// This is a sample recipes API
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Contact: Siddhartha Sitaula <sitaulasiddhartha2002@gmail.com> https://sitaulasiddhartha2002@gmail.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipies []Recipe

func init() {
	recipies = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipies)
}

// swagger:operation POST /recipes recipes NewRecipes
// Creates a new recipe
// ---
//
// requestBody:
//
//	description: Request body of the POST /recipes
//	requried: true
//	content:
//	  application/json
//
// produces:
// - application/json
//
// responses:
//
//	'200':
//		description: Successful Operations
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipies = append(recipies, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes ListRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
//
// responses:
//
//	'200':
//	   description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(200, recipies)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//		description: Successful Operation
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1

	for i := 0; i < len(recipies); i++ {
		if recipies[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Recipe not found",
		})
	}

	recipies[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipies); i++ {
		if recipies[i].ID == id {
			index, _ = strconv.Atoi(id)
		}
	}

	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Recipe does not exist",
		})
		return
	}

	recipies = append(recipies[:index], recipies[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})

}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipies); i++ {
		found := false
		for _, t := range recipies[i].Tags {
			if strings.EqualFold(tag, t) {
				found = true
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipies[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.Run(":8080")
}
