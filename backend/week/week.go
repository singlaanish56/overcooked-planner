package week

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/backend/recipe"
	"github.com/singlaanish56/overcooked-planner/database"
)

type Meal struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
	RecipesID []string `json:"recipesIds"`
	Recipes [] recipe.Recipe `json:"recipes"`
}

type MealDay struct{
	ID string `json:"id"`
	Day string `json:"day"`
	Date string `json:"date"`
	TotalMeals []Meal `json:"meals"`
}

func createMealDay(c *gin.Context){
	var newDay MealDay

	if err:=c.BindJSON(&newDay); err!=nil{
		return 
	}

	tx, err := database.DB.Begin()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"cannot start a transaction"})
		return 
	}

	defer tx.Rollback()

	queryInsertMealDay := `insert into mealDay (day, date) values ($1, $2) returning id`
	var mealDayId string

	err1 := tx.QueryRow(queryInsertMealDay, newDay.Day, newDay.Date).Scan(&mealDayId)
	if err1 != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"failed to insert the mealDay"})
	}

	newDay.ID = mealDayId

	queryInsertMeal := `insert into meal (name, time) values ($1, $2) returning id`
	queryInsertRecipeID := `insert into meal_recipe (recipe_id, meal_id) values ($1, $2)`
	queryInsertMealId := `insert into mealDay_meal (meal_id, mealDay_id) values ($1, $2)`

	for _, meal := range newDay.TotalMeals{
		var mealId string
		err1 = tx.QueryRow(queryInsertMeal, meal.Name, meal.Time).Scan(&mealId)
		if err1 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error with creating the meal"})
			return 
		}

		meal.ID = mealId
		for _, recipeID := range meal.RecipesID{
			_, err1 = tx.Exec(queryInsertRecipeID, recipeID, mealId)
			if err1 != nil{
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Unable to map the recipe to the meal"})
				return
			}
		}

		_, err1 = tx.Exec(queryInsertMealId, mealId, mealDayId)
		if err1 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt not link the meal to the mealday"})
			return
		}
	}

	if err1 = tx.Commit(); err1 != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error during committing the txn"})
		return 
	}
	c.IndentedJSON(http.StatusCreated, newDay)
}

func getAllMealDays(c *gin.Context){

	//c.IndentedJSON(http.StatusOK, days)
}

func getMealDayById(c *gin.Context){
	id := c.Param("id")

	tx, err := database.DB.Begin();
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt start the txn"})
		return
	}

	defer tx.Rollback()

	queryGetMealDay := `select md.id AS mealday_id, md.day, md.date, m.id as meal_id, m.name, m.time,
						r.id as recipe_id, r.name as recipe_name, r.description as recipe_desc, r.people as recipe_people,
	           			i.name as indgredient_name, i.type as indgredient_type,
			 			ri.quantity, ri.unit
						from mealday md
						join mealday_meal mdm on md.id=mdm.mealday_id
						join meal m on mdm.meal_id = m.id
						left join meal_recipe mr on m.id = mr.meal_id
						left join recipe r on mr.recipe_id = r.id
						left join recipe_indgredient ri on r.id =  ri.recipe_id
						left join indgredient i on ri.indgredient_id = i.id
						where md.id = $1
						ORDER BY 
    					m.id, r.id, i.name
						`

	
	rows, err := tx.Query(queryGetMealDay, id)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt retrieve"})
		return
	}

	defer rows.Close()

	var mealDay *MealDay
	mealMap := make(map[string]*Meal)
	recipeMap := make(map[string]*recipe.Recipe)

	for rows.Next(){
		var md MealDay
		var m Meal
		var r recipe.Recipe
		var i recipe.Indgredient
		err:= rows.Scan(&md.ID, &md.Day, &md.Date, 
						&m.ID,&m.Name,&m.Time,
					    &r.ID,&r.Name,&r.Description,&r.People,
						&i.Name,&i.Type,&i.Quantity,&i.Unit)
		if err !=nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error while retrieving the meal day"})
			return
		}

		if mealDay == nil{
			mealDay = &md
			mealDay.TotalMeals = make([]Meal,0)
		}
		
		meal, exists := mealMap[m.ID]
		if !exists{
			meal = &m
			meal.RecipesID = make([]string, 0)
			meal.Recipes = make([]recipe.Recipe, 0)
			mealMap[m.ID] =  meal
			mealDay.TotalMeals = append(mealDay.TotalMeals, *meal)
		}

		rec, exists := recipeMap[r.ID]
		if !exists{
			rec = &r
			rec.Indgredients = make([]recipe.Indgredient, 0)
			recipeMap[r.ID] = rec
			meal.RecipesID = append(meal.RecipesID, r.ID)
			meal.Recipes = append(meal.Recipes, *rec)
		}

		recipeIndex := -1
		for idx, existingRecipe := range meal.Recipes{
			if existingRecipe.ID == rec.ID{
				recipeIndex = idx
				break;
			}
		}

		if recipeIndex != -1{
			meal.Recipes[recipeIndex].Indgredients = append(meal.Recipes[recipeIndex].Indgredients, i)
		}else{
			rec.Indgredients = append(rec.Indgredients, i)
			meal.RecipesID = append(meal.RecipesID, rec.ID)
			meal.Recipes = append(meal.Recipes, *rec)
		}
		
		mealMap[m.ID] = meal

		for i, totalMeal := range mealDay.TotalMeals{
			if totalMeal.ID == meal.ID{
				mealDay.TotalMeals[i] =*meal
			}
		}
	}

	if err := tx.Commit(); err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error while getting meal day information"})
		return
	}

	if mealDay == nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"error":"No Recipe found"})
		return 
	}

	c.IndentedJSON(http.StatusFound, mealDay)
}

func updateMealDay(c *gin.Context){
	id:=c.Param("id")
	
	var updateMealDay MealDay

	if err:=c.BindJSON(&updateMealDay); err!=nil{
		return 
	}


	tx, err := database.DB.Begin();
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt start the txn"})
		return
	}

	defer tx.Rollback()

	
	//retrieve the mealDay and get all the corresponding mealIds in the mapping table
	queryIds := "Select mdm.meal_id from mealday_meal mdm where mdm.mealday_id = $1"
	rows, err := tx.Query(queryIds, id)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error while selecting"})
		return 
	}
	defer rows.Close()
	
	mealIDsInDB := make(map[string]int, 0)
	for rows.Next(){
		var mealId string
		err := rows.Scan(&mealId)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error scanning the row to get the meal IDs"})
			return
		}
		mealIDsInDB[mealId] =1
	}

	queryDeleteTheMealRecipe := "delete from meal_recipe where meal_id=$1"
	queryUpdateTheMeal := "update meal set name=$1, time=$2 where id=$3"
	queryAddTheRecipeId := "insert into meal_recipe (meal_id, recipe_id) values ($1, $2)"
	queryInsertNewMeal := "insert into meal (name, time) values ($1,$2) returning id"
	queryInsertTheMealDayMap := "insert into mealday_meal (mealday_id, meal_id) values($1, $2)"
	for _, meal := range updateMealDay.TotalMeals{
		_, ok := mealIDsInDB[meal.ID]
		if ok {
			_, err := tx.Exec(queryDeleteTheMealRecipe, meal.ID)
			if err != nil{
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error while deleting the recipe row"})
				return
			}
			_, err = tx.Exec(queryUpdateTheMeal, meal.Name, meal.Time, meal.ID)
			if err != nil{
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":err})
				return				
			}

			for _, recipeId := range meal.RecipesID{
				_, err = tx.Exec(queryAddTheRecipeId, meal.ID, recipeId)
				if err != nil{
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error while updating the recipe map"})
					return
				}
			}
			
			delete(mealIDsInDB, meal.ID)
		}else{
			var newmealId string
			err := tx.QueryRow(queryInsertNewMeal, meal.Name, meal.Time).Scan(&newmealId)
			if err != nil{
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error while inserting the new meal"})
				return
			}

			for _, recipeId := range meal.RecipesID{
				_, err = tx.Exec(queryAddTheRecipeId, newmealId, recipeId)
				if err != nil{
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error while updating the recipe map"})
					return
				}
			}

			_, err = tx.Exec(queryInsertTheMealDayMap, id, newmealId)
			if err !=  nil{
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"while updating the mealDay map"})
				return
			}

			meal.ID = newmealId
		}
	}

	queryDeleteRecipeId := ` delete from meal_recipe where meal_id=$1`
	queryDeleteTheMealId := `delete from mealday_meal where meal_id=$1`
	queryDeleteTheMeal := `delete from meal where id=$1`
	for mealid, _ := range mealIDsInDB{
			
		_, err := tx.Exec(queryDeleteRecipeId, mealid)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error deleting the associated recipes"})
			return					
		}

		//query delete the map for mealDay / meal

		_, err1 := tx.Exec(queryDeleteTheMealId, mealid)
		if err1 != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error deleting the associated meals"})
		return					
		}

		//query to delete the meal


		_, err = tx.Exec(queryDeleteTheMeal, mealid)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":err})
			return					
		
	}
}

	if err:=tx.Commit(); err!=nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"error while committing the transaction"})
		return 
	}

	c.IndentedJSON(http.StatusAccepted, updateMealDay)
}

func deleteMealById(c * gin.Context){
	id:=c.Param("id")

	tx, err  := database.DB.Begin()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt start the transaction"})
		return 
	}
	defer tx.Rollback()

	queryGetMeals := `select mdm.meal_id from mealDay_meal mdm where mdm.mealDay_id = $1`
	meals := make([]string,0)

	rows, err := tx.Query(queryGetMeals, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error getting the meals for it"})
		return
	}

	for rows.Next(){
		var mealId string
		err1 := rows.Scan(&mealId)
		if err1 != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error getting the meals"})
			return			
		}

		meals = append(meals, mealId)
	}

	queryDeleteRecipeId := ` delete from meal_recipe where meal_id=$1`
	for _, mealId := range meals{
		_, err := tx.Exec(queryDeleteRecipeId, mealId)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error deleting the associated recipes"})
			return					
		}
	}
	
	//query delete the map for mealDay / meal
	queryDeleteTheMealId := `delete from mealDay_meal where mealDay_id=$1`
	_, err1 := tx.Exec(queryDeleteTheMealId, id)
	if err1 != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error deleting the associated meals"})
		return					
	}

	//query to delete the meal
	queryDeleteTheMeal := `delete from meal where id=$1`
	for _, mealID := range meals{
		_, err := tx.Exec(queryDeleteTheMeal, mealID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Error deleting the meals"})
			return					
		}

	}
	
	queryDeleteTheMealDay := `delete from mealDay where id=$1`
	_, err2 := tx.Exec(queryDeleteTheMealDay, id)
	if err2 != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err2})
		return					
	} 

	if err = tx.Commit(); err !=nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Couldnt Commit the transaction"})
		return 				
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message":"meal Day deleted"})
}

func StartWeek(router *gin.Engine) {
	fmt.Println("Starting the week backend")
	
	router.POST("/mealDay", createMealDay)
	router.POST("/mealDay/:id", updateMealDay)
	router.GET("/C",getAllMealDays)
	router.GET("/mealDay/:id", getMealDayById)
	router.DELETE("/mealDay/:id", deleteMealById)

}