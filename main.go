package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"

	"pawrest/sqldb"
	"github.com/gin-gonic/gin"
)

func getBooks(c *gin.Context) {
	books, err := sqldb.GetKsiazki()

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	book, err := sqldb.GetKsiazka(int64(id))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Nie znaleziono książki"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := sqldb.DelKsiazka(int64(id)); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Nie znaleziono książki"})
		return
	}

	c.Status(204)
}

func main() {
	fmt.Println("Starting the server...")
	sqldb.ConnectToDB()

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("api/v1/books", getBooks)
	router.GET("api/v1/books/:id", getBook)

	router.DELETE("api/v1/books/:id", deleteBook)

	router.Run("localhost:8080")
}
