package recipe

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/database"
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

	tx, err := database.DB.Begin()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to start the txn please try again"})
		return
	}

	defer tx.Rollback()

	query := `insert into recipe (name, description, people) values($1, $2, $3) returning id`
	var recipeID string

	err1 := tx.QueryRow(query, newRecipe.Name, newRecipe.Description, newRecipe.People).Scan(&recipeID);
	if(err1 != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to insert the recipe"})
		return 		
	}

	newRecipe.ID = recipeID
	mapIDIndgredients :=make(map[string] Indgredient)

	queryCheck := "select i.id from indgredient i where LOWER(i.name) like LOWER($1) and LOWER(i.type) LIKE LOWER($2)"
	queryInsert := `insert into indgredient (name, type) values($1, $2) returning id`
	//iterate over the indgredients	
	for _, ing := range newRecipe.Indgredients{	
		var indId string
		//check if the indgredient is present in the table
		err := tx.QueryRow(queryCheck , ing.Name, ing.Type).Scan(&indId)
		if(err== sql.ErrNoRows){
			err = tx.QueryRow(queryInsert, ing.Name, ing.Type).Scan(&indId)
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to insert the indgredient"})
				return 						
			}
		}else if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to insert the indgredient"})
			return 					
		}
		mapIDIndgredients[indId] = ing
		
	}
	
	//insert the values in the recipe_indgredient table
	queryMapRecipeIndgredient := `insert into recipe_indgredient (recipe_id, indgredient_id, quantity, unit) values($1, $2, $3, $4)`
	for key, val := range mapIDIndgredients{
		
		_, err := tx.Exec(queryMapRecipeIndgredient, recipeID, key, val.Quantity, val.Unit)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to map the indgredient"})
			return 				
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to Commit"})
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
	id:=c.Param("id")

	tx, err  := database.DB.Begin()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt start the transaction"})
		return 
	}
	defer tx.Rollback()

	queryDelete := "delete from recipe_indgredient where recipe_id=$1"
	_, err2 := tx.Exec(queryDelete, id)
	if err2!=nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt delete the indgredient in the recipe"})
		return 		
	}

	queryDelete = "delete from recipe where id=$1"
	_, err1 := tx.Exec(queryDelete, id)
	if err1!=nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt delete the recipe"})
		return 		
	}



	if err = tx.Commit(); err !=nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt Commit the transaction"})
		return 				
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message":"recipe deleted"})
}

func updateTheRecipeById(c *gin.Context){
	id:=c.Param("id")
	

	var newRecipe Recipe
	if err:=c.BindJSON(&newRecipe); err!=nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"recipe could not be updated"})
		return 
	}

	tx, err := database.DB.Begin()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"transaction could not be started"})
		return 
	}
	defer tx.Rollback()

	_, err = tx.Exec("update recipe set name=$1 , description=$2, people=$3 where id=$4", newRecipe.Name, newRecipe.Description, newRecipe.People, id)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"Update the recipe failed"})
		return 
	}

	mapIDIndgredients :=make(map[string] Indgredient)
	queryCheck := "select i.id from indgredient i where LOWER(i.name) like LOWER($1) and LOWER(i.type) LIKE LOWER($2)"
	queryInsert := `insert into indgredient (name, type) values($1, $2) returning id`

	for _, ing := range newRecipe.Indgredients{	
		var indId string
		//check if the indgredient is present in the table
		err := tx.QueryRow(queryCheck , ing.Name, ing.Type).Scan(&indId)
		if(err== sql.ErrNoRows){
			err = tx.QueryRow(queryInsert, ing.Name, ing.Type).Scan(&indId)
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to insert the indgredient"})
				return 						
			}
		}else if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to insert the indgredient"})
			return 					
		}
		mapIDIndgredients[indId] = ing
		
	}

	queryMapRecipeIndgredient := `update recipe_indgredient set quantity=$1 , unit=$2 where recipe_id=$3 and indgredient_id=$4`
	for key, val := range mapIDIndgredients{
		
		_, err := tx.Exec(queryMapRecipeIndgredient, val.Quantity, val.Unit, id, key)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to map the indgredient"})
			return 				
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to Commit"})
		return 				
	}
	c.IndentedJSON(http.StatusCreated, newRecipe)

}

func StartRecipe(router *gin.Engine) {
	fmt.Println("This is a start to the recipe backend")

	router.GET("/recipes", getRecipes)
	router.POST("/recipes",addRecipe)
	router.GET("/recipes/:id", getRecipeById)
	router.DELETE("/recipes/:id",deleteRecipeById)
	router.POST("/recipes/:id",updateTheRecipeById)

}

