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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
	"tutorial/handlers"

	"github.com/gin-contrib/sessions"
	sRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// swagger:parameters recipes NewRecipes
type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipies []Recipe
var ctx context.Context
var err error
var client *mongo.Client

var recipiesHandler *handlers.RecipeHandler
var authHandler *handlers.AuthHandler

func init() {

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27027/"))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	if err != err {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	collection := client.Database("tutorial").Collection("recipes")
	collectionUsers := client.Database("tutorial").Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	recipiesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
	status := redisClient.Ping()
	fmt.Println(status)

	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
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

	collection := client.Database("tutorial").Collection("recipes")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)

	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
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
	objectID, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database("tutorial").Collection("recipes")

	_, err := collection.UpdateOne(ctx, bson.M{
		"_id": objectID,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated",
	})
}

// func DeleteRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	index := -1

// 	for i := 0; i < len(recipies); i++ {
// 		if recipies[i].ID == id {
// 			index, _ = strconv.Atoi(id)
// 		}
// 	}

// 	if index == -1 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Recipe does not exist",
// 		})
// 		return
// 	}

// 	recipies = append(recipies[:index], recipies[index+1:]...)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Recipe has been deleted",
// 	})

// }

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)

// 	for i := 0; i < len(recipies); i++ {
// 		found := false
// 		for _, t := range recipies[i].Tags {
// 			if strings.EqualFold(tag, t) {
// 				found = true
// 			}
// 		}

// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipies[i])
// 		}
// 	}

// 	c.JSON(http.StatusOK, listOfRecipes)
// }

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
// 			c.AbortWithStatus(401)
// 		}
// 		c.Next()
// 	}
// }

func main() {
	users := map[string]string{
		"admin": "admin",
		"user":  "user",
	}
	h := sha256.New()
	ctx := context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27027/"))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collectionUsers := client.Database("tutorial").Collection("users")
	// store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	for username, password := range users {
		collectionUsers.InsertOne(ctx, bson.M{
			"username": username,
			"password": hex.EncodeToString(h.Sum([]byte(password))),
		})
	}

	router := gin.Default()
	store, _ := sRedis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	authorized := router.Group("/")
	router.GET("/recipes", recipiesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	authorized.Use(authHandler.AuthMiddleware())
	{

		authorized.POST("/recipes", recipiesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", UpdateRecipeHandler)
		authorized.POST("/signout", authHandler.SignOutHandler)
		// router.DELETE("/recipes/:id", DeleteRecipeHandler)
		// router.GET("/recipes/search", SearchRecipesHandler)
	}

	router.Run(":8080")
}
