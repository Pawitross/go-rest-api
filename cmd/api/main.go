package main

import (
	"fmt"
	"log"

	"pawrest/internal/api/routes"
	"pawrest/internal/db"
)

func main() {
	fmt.Println("Łączenie z bazą danych...")
	if err := db.ConnectToDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	fmt.Println("Uruchamianie serwera...")
	//gin.SetMode(gin.ReleaseMode)
	router := routes.Router()

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
