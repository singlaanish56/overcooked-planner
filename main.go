package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/singlaanish56/overcooked-planner/api"
	"github.com/singlaanish56/overcooked-planner/database"
)

func getDbConnectionUrl() string{
	if err := godotenv.Load(); err != nil{
		fmt.Println("Coudlnt load the env variables")
		os.Exit(1)
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("user"), os.Getenv("password"), os.Getenv("host"), os.Getenv("port"), os.Getenv("db_name"))
	fmt.Println(connStr)
	return connStr
}

func main(){
	router := gin.Default()

	conn, err := pgx.Connect(context.Background(), getDbConnectionUrl())
	if err !=nil{
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	queries := database.New(conn)

	api.RegisterProductRoutes(router, queries)

	router.Run("localhost:8080")
}