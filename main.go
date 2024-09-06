package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"strings"
	_ "terrenceng/recipes-api/docs"
	"time"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success string `json:"message"`
}

// NewRecipeHandler godoc
// @Summary      Create a recipe
// @Description  Create a recipe
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  Recipe
// @Failure 	 400  {object}  Error
// @Router       /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler godoc
// @Summary      List recipes
// @Description  Get all recipes
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  []Recipe
// @Router       /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler godoc
// @Summary      Update a recipe
// @Description  Update a recipe
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  Recipe
// @Failure 	 400  {object}  Error
// @Failure 	 404  {object}  Error
// @Router       /recipe/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipeHandler godoc
// @Summary      Delete a recipe
// @Description  Delete a recipe by ID
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  Success
// @Failure 	 404  {object}  Error
// @Router       /recipe/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// SearchRecipesHandler godoc
// @Summary      Search recipes
// @Description  Search recipes by tags
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  []Recipe
// @Router       /recipe/tag [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

// @title           Recipes API
// @version         1.0
// @description     This is a sample recipes API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Terrence NG
// @contact.url
// @contact.email  kh.terence.ng@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:7778
// @BasePath

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	router := gin.Default()
	router.POST("/recipe", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipe/:id", UpdateRecipeHandler)
	router.DELETE("/recipe/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":7778")
}
