package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pawrest/internal/db"
	m "pawrest/internal/models"
)

func validateBook(b m.Ksiazka) bool {
	return b.Tytul == "" ||
		b.Rok == 0 ||
		b.Strony <= 0 ||
		b.Autor <= 0 ||
		b.Gatunek <= 0 ||
		b.Jezyk <= 0
}

func GetBooks(c *gin.Context) {
	params := c.Request.URL.Query()

	books, err := db.GetKsiazki(params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano nieodpowiedni identyfikator"})
		return
	}

	book, err := db.GetKsiazka(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func PostBook(c *gin.Context) {
	var newBook m.Ksiazka

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if validateBook(newBook) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano puste lub nieprawidłowe pola"})
		return
	}

	id, err := db.InsertKsiazka(newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook.Id = id

	location := "/books/" + strconv.FormatInt(newBook.Id, 10)
	c.Header("Location", location)

	c.JSON(http.StatusCreated, newBook)
}

func PutBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano nieodpowiedni identyfikator"})
		return
	}

	var newBook m.Ksiazka

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if validateBook(newBook) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano puste lub nieprawidłowe pola"})
		return
	}

	if err := db.UpdateWholeKsiazka(int64(id), newBook); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func PatchBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano nieodpowiedni identyfikator"})
		return
	}

	var patchBook m.Ksiazka

	if err := c.BindJSON(&patchBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if err := db.UpdateKsiazka(int64(id), patchBook); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano nieodpowiedni identyfikator"})
		return
	}

	if err := db.DelKsiazka(int64(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
