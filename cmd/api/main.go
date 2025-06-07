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
	fmt.Println("Łączenie z bazą danych...")
	if err := db.ConnectToDB(); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	if err := middleware.InitLogger(); err != nil {
		log.Fatalf("Nie udało się zainicjalizować logowania żądań: %v\n", err)
	}
	defer middleware.CloseLogger()

	fmt.Println("Uruchamianie serwera...")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(middleware.FileLogger())
	routes.Router(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
