package api

import (
	"github.com/gin-gonic/gin"
)

func getAll(c *gin.Context){
	
}

func RegisterProductRoutes(router *gin.Engine){
	fmt.Println("Starting the list backend")

	router.GET("/items", getAll)
	router.GET("/items/item/:id", getByItemId)

	router.POST("/items/item", addItem)
	router.POST("/items/item/:id", updateItemById)
	
}
