package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pawrest/internal/db"
	m "pawrest/internal/models"
)

// @Summary		Get a list of all books
// @Description	Responds with a list of all books as JSON. Optional filtering and pagination is available through parameters.
// @Tags			Books
// @Produce		json
// @Param			id			query		string			false	"Book id"
// @Param			title		query		string			false	"Book title"
// @Param			year		query		int				false	"Year of publishing of the book"
// @Param			pages		query		int				false	"Number of pages in the book"
// @Param			author		query		int				false	"Author id"
// @Param			genre		query		int				false	"Genre id"
// @Param			language	query		int				false	"Language id"
// @Param			limit		query		int				false	"Limit returned number of resources"
// @Param			offset		query		int				false	"Offset returned resources"
// @Success		200			{array}		models.Book		"OK - Fetched books"
// @Failure		404			{object}	models.Error	"Not Found"
// @Router			/books [get]
func GetBooks(c *gin.Context) {
	params := c.Request.URL.Query()

	books, err := db.GetBooks(params)
	if err != nil {
		c.JSON(http.StatusNotFound, m.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary		Get one book
// @Description	Responds with the queried book as JSON or an error message.
// @Tags			Books
// @Produce		json
// @Param			id	path		int				true	"Book id"
// @Success		200	{object}	models.Book		"OK - Fetched book"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid book id"
// @Failure		404	{object}	models.Error	"Not Found"
// @Router			/books/{id} [get]
func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Provided incorrect identifier"})
		return
	}

	book, err := db.GetBook(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, m.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary		Create a new book
// @Description	Accepts a JSON body to create a new book. Responds with the created book and set Location header or an error message.
// @Tags			Books
// @Accept			json
// @Produce		json
// @Param			book	body		models.Book		true	"New Book"
// @Success		201		{object}	models.Book		"Created - Added new book"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Header			201		{string}	Location		"Path of the newly created book"
// @Router			/books [post]
func PostBook(c *gin.Context) {
	var newBook m.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newBook.ValidateBook() {
		c.JSON(http.StatusBadRequest, m.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	id, err := db.InsertBook(newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: err.Error()})
		return
	}

	newBook.Id = id

	location := c.FullPath() + "/" + strconv.FormatInt(newBook.Id, 10)
	c.Header("Location", location)

	c.JSON(http.StatusCreated, newBook)
}

// @Summary		Update an existing book
// @Description	Accepts a JSON body to update a book. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Books
// @Accept			json
// @Param			id		path	int			true	"Existing Book id"
// @Param			book	body	models.Book	true	"Updated Book"
// @Success		204		"Updated the book"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Router			/books/{id} [put]
func PutBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Provided incorrect identifier"})
		return
	}

	var newBook m.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newBook.ValidateBook() {
		c.JSON(http.StatusBadRequest, m.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	if err := db.UpdateWholeBook(int64(id), newBook); err != nil {
		c.JSON(http.StatusNotFound, m.Error{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Patch an existing book
// @Description	Accepts a JSON body with patch data to a book. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Books
// @Accept			json
// @Param			id		path	int			true	"Existing Book id"
// @Param			book	body	models.Book	true	"Patches to the book"
// @Success		204		"No Content - Successfully patched the book"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Router			/books/{id} [patch]
func PatchBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Provided incorrect identifier"})
		return
	}

	var patchBook m.Book

	if err := c.BindJSON(&patchBook); err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Invalid JSON in request body"})
		return
	}

	if err := db.UpdateBook(int64(id), patchBook); err != nil {
		c.JSON(http.StatusNotFound, m.Error{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Delete an existing book
// @Description	Responds with a status code. When an error occurs the response body contains an error message.
// @Tags			Books
// @Param			id	path	int	true	"Book id"
// @Success		204	"No Content - Successfully deleted the book"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid book id"
// @Router			/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.Error{Error: "Provided incorrect identifier"})
		return
	}

	if err := db.DelBook(int64(id)); err != nil {
		c.JSON(http.StatusNotFound, m.Error{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
