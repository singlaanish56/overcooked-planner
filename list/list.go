package list

import (
	"fmt"
	"github.com/singlaanish56/overcooked-planner/recipe"
	"github.com/gin-gonic/gin"
	"net/http"
)
type ListItem struct{
	Id 			string `json:"id"`
	Indgredient recipe.Indgredient `json:"indgredient"`
	IsAvailable bool `json:"isAvailable"`
} 

type List struct{
	Id    string `json:"id"`
	IndgredientList []ListItem `json:"indgredientList"` 
}

var list = []List{
	{Id:"1",IndgredientList: []ListItem{
		{Id:"1", Indgredient: recipe.Indgredient{Name:"Chicken",Quantity: 900,Unit: "g"}, IsAvailable: false},
		{Id:"2",Indgredient: recipe.Indgredient{Name:"Wheat Parantha",Quantity: 12,Unit:"pc"},IsAvailable: false},
		{Id:"3",Indgredient: recipe.Indgredient{Name:"Dal",Quantity: 200,Unit:"g"},IsAvailable: false},
		},
	},	
}

func GetAllList(c *gin.Context){
	c.IndentedJSON(http.StatusOK, list)
}

func GetListItemById(c *gin.Context){
	id := c.Param("id")

	for _,a := range list{
		if a.Id == id{
			c.IndentedJSON(http.StatusFound, a)
			return
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"item not found"})
	}
}

func GetListItemByName(c *gin.Context){
	name := c.Param("name")

	for _,a := range list{
		for _,b := range a.IndgredientList{
			if b.Indgredient.Name == name{
				
				c.IndentedJSON(http.StatusOK, b)
				return 
			}
		}
	}
}

func UpdateListById(c *gin.Context){
	id := c.Param("id")

	for i,a := range list{
		if a.Id == id{
			var item List
			if err:= c.BindJSON(&item); err!=nil{
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"item could not be updated"})
				return 
			}
			list[i] = item
			c.IndentedJSON(http.StatusCreated, gin.H{"message":"item updated"})
			return 
		}
	}
}

func UpdateItemByName(c *gin.Context){

	var updateItem recipe.Indgredient
	if err:= c.BindJSON(&updateItem); err !=nil{
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"the item could not be updated"})
			return
	}

	for _,a := range list{
		for j,b := range a.IndgredientList{
			if b.Indgredient.Name == updateItem.Name{
				a.IndgredientList[j].Indgredient = updateItem
				c.IndentedJSON(http.StatusCreated, gin.H{"message":"the item is updated"})
				return
			}
		}
	}
}


func StartListBackend() {
	fmt.Println("Start the List backend")
	router := gin.Default()

	router.GET("/lists",GetAllList)
	router.GET("/lists/:id", GetListItemById)
	router.GET("/lists/name/:name", GetListItemByName)
	router.POST("/lists/:id", UpdateListById)
	router.POST("/lists", UpdateItemByName)
	
	router.Run("localhost:8080")
}