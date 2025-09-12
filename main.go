package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/singlaanish56/overcooked-planner/api"
	"github.com/singlaanish56/overcooked-planner/database"
)


func main(){
	router := gin.Default()

	conn := database.ConnectionDb()

	defer conn.Close(context.Background())

	queries := database.New(conn)

	insertedProduct, err := queries.CreateProduct(context.Background(), database.CreateProductParams{Name: "lactose free milk", Company: pgtype.Text{String:"amul"}, Subtype: pgtype.Text{String: "consumable"}, Weight: pgtype.Float4{Float32: 250}, Unit: pgtype.Text{String:"ml"}})
	if err !=nil{
		fmt.Sprintf("Couldn't do the database transaction %s", err)
	}
	log.Println(insertedProduct)
	fetchedProduct, err := queries.GetProduct(context.Background(), insertedProduct.Name)
	if err !=nil{
		fmt.Sprintf("Couldn't do the database transaction 2 %s", err)
	}
	
	log.Println(fetchedProduct)
	api.RegisterProductRoutes(router)

	router.Run("localhost:8080")
}