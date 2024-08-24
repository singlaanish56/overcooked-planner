package week

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/backend/recipe"
)

type Meal struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
	RecipesID []string `json:"recipes"`
	Recipes [] recipe.Recipe
}

type MealDay struct{
	ID string `json:"id"`
	Day string `json:"day"`
	Date string `json:"date"`
	TotalMeals []Meal `json:"meals"`
}

//seed for the test
var recipes = []recipe.Recipe{
	{ID:"1",Name: "Chicken Roll",People:2,Indgredients: []recipe.Indgredient{
		{Name:"Chicken",Quantity: 200,Unit: "g"},
		{Name:"Wheat Parantha",Quantity: 4,Unit:"pc"},
		{Name:"Marinate",Quantity: 30,Unit: "g"},
	}},
	{ID:"2",Name: "Khichdi",People:2,Indgredients: []recipe.Indgredient{
		{Name:"Rice",Quantity: 100,Unit: "g"},
		{Name:"Dal",Quantity: 100,Unit:"g"},
	}},
	{ID:"3",Name: "Half Egg",People:1,Indgredients: []recipe.Indgredient{
		{Name:"Egg",Quantity: 2,Unit: "pc"},
		{Name:"Bread",Quantity: 2,Unit:"pc"},
	}},
	}


var meals = []Meal{
{ID:"1",Name:"Breakfast",Time:"Morning",RecipesID: []string{"1"},Recipes: []recipe.Recipe{recipes[0]}},
{ID:"2",Name:"Lunch",Time:"Afternoon",RecipesID: []string{"2"},Recipes: []recipe.Recipe{recipes[1]}},
{ID:"3",Name:"Dinner",Time:"Night",RecipesID: []string{"3"},Recipes: []recipe.Recipe{recipes[2]}},
}

var days = []MealDay{
{ID:"1",Day:"Monday",Date:"11-06-24",TotalMeals: meals},
{ID:"2",Day:"Tuesday",Date:"12-06-24",TotalMeals: meals},
{ID:"3",Day:"Wednesday",Date:"13-06-24",TotalMeals: meals},
}

func createMealDay(c *gin.Context){
	var newDay MealDay

	if err:=c.BindJSON(&newDay); err!=nil{
		return 
	}

	for _, a := range newDay.TotalMeals{
		for _,b :=range a.RecipesID{
			for _,c := range recipes{
				if(b==c.ID){
					a.Recipes = append(a.Recipes, c)
				}
			}
			
		}
	}

	days = append(days, newDay)
	c.IndentedJSON(http.StatusCreated, newDay)
}

func getAllMealDays(c *gin.Context){

	c.IndentedJSON(http.StatusOK, days)
}

func getMealDayById(c *gin.Context){
	id := c.Param("id")

	for _, a := range days{
		if(a.ID == id){
			c.IndentedJSON(http.StatusFound, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"week not found"})
}

func updateMealDay(c *gin.Context){
	id:=c.Param("id")
	
	for i,a := range days{
		if a.ID == id{
			var newMeal MealDay
			if err:=c.BindJSON(&newMeal); err!=nil{
				fmt.Println(err)
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"meal could not be updated"})
				return 
			}


			for _, a := range newMeal.TotalMeals{
				for _,b :=range a.RecipesID{
					for _,c := range recipes{
						if(b==c.ID){
							a.Recipes = append(a.Recipes, c)
						}
					}
					
				}
			}

			days[i]=newMeal

			c.IndentedJSON(http.StatusCreated, gin.H{"message":"meal updated"})
			return 
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"recipe not found"})

}

func deleteMealById(c * gin.Context){
	id:=c.Param("id")
	for i,a := range days{
		if a.ID == id{
			days = append(days[:i],days[i+1:]...)
			return 
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"meal not found"})
}

func StartWeek(router *gin.Engine) {
	fmt.Println("Starting the week backend")
	
	router.POST("/mealDay", createMealDay)
	router.POST("/mealDay/:id", updateMealDay)
	router.GET("/mealDay",getAllMealDays)
	router.GET("/mealDay/:id", getMealDayById)
	router.DELETE("/mealDay/:id", deleteMealById)

}