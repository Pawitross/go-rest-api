package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pawrest/sqldb"
)

func getBooks(c *gin.Context) {
	params := c.Request.URL.Query()

	books, err := sqldb.GetKsiazki(params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	book, err := sqldb.GetKsiazka(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func postBook(c *gin.Context) {
	var newBook sqldb.Ksiazka

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	id, err := sqldb.InsertKsiazka(newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook.Id = id

	c.JSON(http.StatusCreated, newBook)
}

func putBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var newBook sqldb.Ksiazka

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if err := sqldb.UpdateWholeKsiazka(int64(id), newBook); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := sqldb.DelKsiazka(int64(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func main() {
	fmt.Println("Łączenie z bazą danych...")
	if err := sqldb.ConnectToDB(); err != nil {
		log.Fatal(err)
	}
	defer sqldb.Db.Close()

	fmt.Println("Uruchamianie serwera...")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("books", getBooks)
			v1.GET("books/:id", getBook)

			v1.POST("books", postBook)

			v1.PUT("books/:id", putBook)

			v1.DELETE("books/:id", deleteBook)
		}
	}

	router.Run(":8080")
}
