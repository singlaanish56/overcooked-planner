package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/database"
)

var query *database.Queries

func RegisterProductRoutes(router *gin.Engine, queries *database.Queries) {
	fmt.Println("Starting the list backend")

	query = queries
	router.GET("/items", getAll)
	router.GET("/items/item/:name", getByItemName)

	router.POST("/items/item", addItem)

	router.PUT("/items/item/:name", updateItemByName)
}

func getAll(c *gin.Context) {

	products, err := query.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	fmt.Println(products)
	c.JSON(http.StatusOK, products)
}

func getByItemName(c *gin.Context) {
	id := c.Param("name")

	//query call to database
	//implement txn;

	product, err := query.GetProduct(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func addItem(c *gin.Context) {

	var createProductParams database.CreateProductParams
	if err := c.BindJSON(&createProductParams); err != nil {
		return
	}

	//query call to the database
	//implement the txn

	product, err := query.CreateProduct(context.Background(), createProductParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func updateItemByName(c *gin.Context) {
	id := c.Param("name")

	product, err := query.GetProduct(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	var updatesProductParams database.UpdateProductParams

	if err := c.BindJSON(&updatesProductParams); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	updatesProductParams.ID = product.ID

	err = query.UpdateProduct(context.Background(), updatesProductParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error updating the product in the database %s", err)})
		return
	}

	if len(updatesProductParams.Name) == 0 {
		updatesProductParams.Name = product.Name
	}

	c.JSON(http.StatusOK, updatesProductParams)
}
