package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pawrest/internal/db"
	m "pawrest/internal/models"
)

//	@Summary		Get list of books in an array
//	@Description	Responds with the list of all books as JSON.
//	@Tags			Books
//	@Produce		json
//	@Success		200	{array}	models.Book
//	@Router			/books [get]
func GetBooks(c *gin.Context) {
	params := c.Request.URL.Query()

	books, err := db.GetBooks(params)
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

	book, err := db.GetBook(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func PostBook(c *gin.Context) {
	var newBook m.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if newBook.ValidateBook() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano puste lub nieprawidłowe pola"})
		return
	}

	id, err := db.InsertBook(newBook)
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

	var newBook m.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if newBook.ValidateBook() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano puste lub nieprawidłowe pola"})
		return
	}

	if err := db.UpdateWholeBook(int64(id), newBook); err != nil {
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

	var patchBook m.Book

	if err := c.BindJSON(&patchBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wystąpił problem z JSON"})
		return
	}

	if err := db.UpdateBook(int64(id), patchBook); err != nil {
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

	if err := db.DelBook(int64(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
