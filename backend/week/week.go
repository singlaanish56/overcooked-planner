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
	Recipes [] recipe.Recipe `json:"recipes"`
}

type MealDay struct{
	ID string `json:"id"`
	Day string `json:"day"`
	Date string `json:"date"`
	TotalMeals []Meal `json:"meals"`
}

//seed for the test
// var recipes = []recipe.Recipe{
// 	{ID:"1",Name: "Chicken Roll",People:2,Indgredients: []recipe.Indgredient{
// 		{Name:"Chicken",Quantity: 200,Unit: "g"},
// 		{Name:"Wheat Parantha",Quantity: 4,Unit:"pc"},
// 		{Name:"Marinate",Quantity: 30,Unit: "g"},
// 	}},
// 	{ID:"2",Name: "Khichdi",People:2,Indgredients: []recipe.Indgredient{
// 		{Name:"Rice",Quantity: 100,Unit: "g"},
// 		{Name:"Dal",Quantity: 100,Unit:"g"},
// 	}},
// 	{ID:"3",Name: "Half Egg",People:1,Indgredients: []recipe.Indgredient{
// 		{Name:"Egg",Quantity: 2,Unit: "pc"},
// 		{Name:"Bread",Quantity: 2,Unit:"pc"},
// 	}},
// 	}


// var meals = []Meal{
// {ID:"1",Name:"Breakfast",Time:"Morning",RecipesID: []string{"1"},Recipes: []recipe.Recipe{recipes[0]}},
// {ID:"2",Name:"Lunch",Time:"Afternoon",RecipesID: []string{"2"},Recipes: []recipe.Recipe{recipes[1]}},
// {ID:"3",Name:"Dinner",Time:"Night",RecipesID: []string{"3"},Recipes: []recipe.Recipe{recipes[2]}},
// }

// var days = []MealDay{
// {ID:"1",Day:"Monday",Date:"11-06-24",TotalMeals: meals},
// {ID:"2",Day:"Tuesday",Date:"12-06-24",TotalMeals: meals},
// {ID:"3",Day:"Wednesday",Date:"13-06-24",TotalMeals: meals},
// }

func createMealDay(c *gin.Context){
	// var newDay MealDay

	// if err:=c.BindJSON(&newDay); err!=nil{
	// 	return 
	// }

	// for _, a := range newDay.TotalMeals{
	// 	for _,b :=range a.RecipesID{
	// 		for _,c := range recipes{
	// 			if(b==c.ID){
	// 				a.Recipes = append(a.Recipes, c)
	// 			}
	// 		}
			
	// 	}
	// }

	// days = append(days, newDay)
	// c.IndentedJSON(http.StatusCreated, newDay)
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
			meal.Recipes = make([]recipe.Recipe, 0)
			mealMap[m.ID] =  meal
			mealDay.TotalMeals = append(mealDay.TotalMeals, *meal)
		}

		rec, exists := recipeMap[r.ID]
		if !exists{
			rec = &r
			rec.Indgredients = make([]recipe.Indgredient, 0)
			recipeMap[r.ID] = rec
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

	c.IndentedJSON(http.StatusFound, mealDay)
}

func updateMealDay(c *gin.Context){
	// id:=c.Param("id")
	
	// for i,a := range days{
	// 	if a.ID == id{
	// 		var newMeal MealDay
	// 		if err:=c.BindJSON(&newMeal); err!=nil{
	// 			fmt.Println(err)
	// 			c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"meal could not be updated"})
	// 			return 
	// 		}


	// 		for _, a := range newMeal.TotalMeals{
	// 			for _,b :=range a.RecipesID{
	// 				for _,c := range recipes{
	// 					if(b==c.ID){
	// 						a.Recipes = append(a.Recipes, c)
	// 					}
	// 				}
					
	// 			}
	// 		}

	// 		days[i]=newMeal

	// 		c.IndentedJSON(http.StatusCreated, gin.H{"message":"meal updated"})
	// 		return 
	// 	}
	// }

	// c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})

}

func deleteMealById(c * gin.Context){
	// id:=c.Param("id")
	// for i,a := range days{
	// 	if a.ID == id{
	// 		days = append(days[:i],days[i+1:]...)
	// 		return 
	// 	}
	// }
	// c.IndentedJSON(http.StatusNotFound, gin.H{"message":"meal not found"})
}

func StartWeek(router *gin.Engine) {
	fmt.Println("Starting the week backend")
	
	router.POST("/mealDay", createMealDay)
	router.POST("/mealDay/:id", updateMealDay)
	router.GET("/C",getAllMealDays)
	router.GET("/mealDay/:id", getMealDayById)
	router.DELETE("/mealDay/:id", deleteMealById)

}