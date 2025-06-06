package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"pawrest/internal/api/handler"
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
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("books", handler.GetBooks)
			v1.GET("books/:id", handler.GetBook)

			v1.POST("books", handler.PostBook)

			v1.PUT("books/:id", handler.PutBook)

			v1.DELETE("books/:id", handler.DeleteBook)
		}
	}

	router.Run(":8080")
}
