package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"pawrest/internal/api/middleware"
	"pawrest/internal/api/routes"
	"pawrest/internal/db"
)

func main() {
	fmt.Println("Connecting to the database...")
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	if err := middleware.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logging middleware: %v\n", err)
	}
	defer middleware.CloseLogger()

	fmt.Println("Starting up the server...")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(middleware.FileLogger())
	routes.Router(router, database)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
