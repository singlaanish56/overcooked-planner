package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/database"
)

var query *database.Queries

func RegisterProductRoutes(router *gin.Engine, queries *database.Queries){
	fmt.Println("Starting the list backend")

	query = queries
	router.GET("/items", getAll)
	router.GET("/items/item/:name", getByItemId)

	router.POST("/items/item", addItem)
	router.POST("/items/item/:id", updateItemById)
	
}


func getAll(c *gin.Context){

}

func getByItemId(c *gin.Context){
	id := c.Param("name")

	//query call to database
	//implement txn;

	product, err := query.GetProduct(context.Background(), id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if product ==nil{
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func addItem(c *gin.Context){


	var createProductParams database.CreateProductParams
	if err := c.BindJSON(&createProductParams); err != nil{
		return 
	}	

	//query call to the database
	//implement the txn

	product, err := query.CreateProduct(context.Background(), createProductParams)
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return 
	}

	c.JSON(http.StatusOK, product)
}

func updateItemById(c *gin.Context){

}
