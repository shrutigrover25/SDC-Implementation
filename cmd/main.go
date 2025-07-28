package main

import (
	"mercor/internal/db"
	router "mercor/internal/domain/router"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.Connect()
	db.Seed(database)

	r := gin.Default()
	router.InitRoutes(r)
	r.Run(":8080")
}
