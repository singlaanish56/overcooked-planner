package main

import (
	"net/http"
	"github.com/singlaanish56/overcooked-planner/api"
	"github.com/gin-gonic/gin"
)


func main(){
	router := gin.Default()

	api.RegisterProductRoutes(router)
}