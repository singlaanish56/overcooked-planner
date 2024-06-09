package recipe

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Indgredient struct{
Name string `json:"name"`
Quantity float64 `json:"quantity"`
Unit string `json:"unit"`
}

type Recipe struct{
	ID string `json:"id"`
	Name string `json:"name"`
	People int `json:"people"`
	Indgredients []Indgredient `json:"indgredients"`
}

//seed for the test
var recipes = []Recipe{
{ID:"1",Name: "Chicken Roll",People:2,Indgredients: []Indgredient{
	{Name:"Chicken",Quantity: 200,Unit: "g"},
	{Name:"Wheat Parantha",Quantity: 4,Unit:"pc"},
	{Name:"Marinate",Quantity: 30,Unit: "g"},
}},
{ID:"2",Name: "Khichdi",People:2,Indgredients: []Indgredient{
	{Name:"Rice",Quantity: 100,Unit: "g"},
	{Name:"Dal",Quantity: 100,Unit:"g"},
}},
}

func getRecipes(c *gin.Context){
	c.IndentedJSON(http.StatusOK, recipes)
}

func addRecipe(c *gin.Context){
	var newRecipe Recipe

	if err:=c.BindJSON(&newRecipe); err!=nil{
		return 
	}

	recipes= append(recipes, newRecipe)
	c.IndentedJSON(http.StatusCreated, newRecipe)
}

func getRecipeById(c *gin.Context){
	id:=c.Param("id")

	for _,a := range recipes{
		if a.ID == id{
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})
}

func deleteRecipeById(c * gin.Context){
	id:=c.Param("id")
	for i,a := range recipes{
		if a.ID == id{
			recipes = append(recipes[:i],recipes[i+1:]...)
			return 
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})
}

func updateTheRecipeById(c *gin.Context){
	id:=c.Param("id")
	
	for i,a := range recipes{
		if a.ID == id{
			var newRecipe Recipe
			if err:=c.BindJSON(&newRecipe); err!=nil{
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"recipe could not be updated"})
				return 
			}
			recipes[i]=newRecipe
			c.IndentedJSON(http.StatusCreated, gin.H{"message":"recipe updated"})
			return 
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})

}

func StartRecipe() {
	fmt.Println("This is a start to the recipe backend")
	router := gin.Default()
	router.GET("/recipes", getRecipes)
	router.POST("/recipes",addRecipe)
	router.GET("/recipes/:id", getRecipeById)
	router.DELETE("/recipes/:id",deleteRecipeById)
	router.POST("/recipes/:id",updateTheRecipeById)
	router.Run("localhost:8080")
}

