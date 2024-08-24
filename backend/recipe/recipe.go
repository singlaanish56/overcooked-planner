package recipe

import (
	"fmt"
	"net/http"
	"github.com/singlaanish56/overcooked-planner/database"
	"github.com/gin-gonic/gin"
)

type Indgredient struct{
Name string `json:"name"`
Quantity float64 `json:"quantity"`
Unit string `json:"unit"`
Type string `json:"type"`
}

type Recipe struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	People int `json:"people"`
	Indgredients []Indgredient `json:"indgredients"`
}


//seed for the test
// var recipes = []Recipe{
// {ID:"1",Name: "Chicken Roll",People:2,Indgredients: []Indgredient{
// 	{Name:"Chicken",Quantity: 200,Unit: "g"},
// 	{Name:"Wheat Parantha",Quantity: 4,Unit:"pc"},
// 	{Name:"Marinate",Quantity: 30,Unit: "g"},
// }},
// {ID:"2",Name: "Khichdi",People:2,Indgredients: []Indgredient{
// 	{Name:"Rice",Quantity: 100,Unit: "g"},
// 	{Name:"Dal",Quantity: 100,Unit:"g"},
// }},
// }

// to be updated
func getRecipes(c *gin.Context){
	
	//rows, err := database.DB.Query(SELECT)
}

func addRecipe(c *gin.Context){
	var newRecipe Recipe


	if err:=c.BindJSON(&newRecipe); err!=nil{
		return 
	}


	
	c.IndentedJSON(http.StatusCreated, newRecipe)
}

func getRecipeById(c *gin.Context){
	id:=c.Param("id")
	query := `Select r.id as recipe_id, r.name as recipe_name, r.description as recipe_desc, r.people as recipe_people,
	           i.name as indgredient_name, i.type as indgredient_type,
			  ri.quantity, ri.unit
			  from recipe r
			  join recipe_indgredient ri on r.id = ri.recipe_id
			  join indgredient i on ri.indgredient_id = i.id
			  where r.id = $1
			  order by i.name`

	rows, err := database.DB.Query(query, id)
	if(err!=nil){
		fmt.Println("Read Query For Recipe failed", err)
		return 
	}
	defer rows.Close()

	var recipe *Recipe
	var indgredients []Indgredient


	for rows.Next(){
		var sr Recipe
		var ing Indgredient
		err:=rows.Scan(&sr.ID, &sr.Name, &sr.Description, &sr.People, &ing.Name, &ing.Type, &ing.Quantity, &ing.Unit)
		if err!=nil{
			fmt.Println("Couldn't Convert to Json", err)
			return
		}
		
		indgredients = append(indgredients, ing)
		if(recipe == nil){
			recipe = &sr
		}
	}
	
	if recipe!=nil{
		recipe.Indgredients = append(recipe.Indgredients, indgredients...)
		c.IndentedJSON(http.StatusFound, recipe)
	}else{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})
	}
	
}

func deleteRecipeById(c * gin.Context){
	// id:=c.Param("id")
	// for i,a := range recipes{
	// 	if a.ID == id{
	// 		recipes = append(recipes[:i],recipes[i+1:]...)
	// 		return 
	// 	}
	// }
	// c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})
}

func updateTheRecipeById(c *gin.Context){
	// id:=c.Param("id")
	
	// for i,a := range recipes{
	// 	if a.ID == id{
	// 		var newRecipe Recipe
	// 		if err:=c.BindJSON(&newRecipe); err!=nil{
	// 			c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"recipe could not be updated"})
	// 			return 
	// 		}
	// 		recipes[i]=newRecipe
	// 		c.IndentedJSON(http.StatusCreated, gin.H{"message":"recipe updated"})
	// 		return 
	// 	}
	// }

	// c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})

}

func StartRecipe(router *gin.Engine) {
	fmt.Println("This is a start to the recipe backend")

	router.GET("/recipes", getRecipes)
	router.POST("/recipes",addRecipe)
	router.GET("/recipes/:id", getRecipeById)
	router.DELETE("/recipes/:id",deleteRecipeById)
	router.POST("/recipes/:id",updateTheRecipeById)

}

