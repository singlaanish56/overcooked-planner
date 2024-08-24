package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/singlaanish56/overcooked-planner/backend/list"
	"github.com/singlaanish56/overcooked-planner/backend/recipe"
	"github.com/singlaanish56/overcooked-planner/backend/week"
	"github.com/singlaanish56/overcooked-planner/database"
)

func main() {

	fmt.Println("Starting the backend")
	router := gin.Default()
	
	recipe.StartRecipe(router)
	week.StartWeek(router)
	list.StartList(router)

	database.ConnectionDb()

	router.Run("localhost:8080")
}