package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

// @Summary		Get a list of all books
// @Description	Responds with a list of all books as JSON. Optional filtering, sorting and pagination is available through parameters.
// @Tags			Books
// @Produce		json
// @Param			id					query		string			false	"Book id"
// @Param			title				query		string			false	"Book title"
// @Param			year				query		int				false	"Year of publishing of the book"
// @Param			pages				query		int				false	"Number of pages in the book"
// @Param			author				query		int				false	"Author id"
// @Param			genre				query		int				false	"Genre id"
// @Param			language			query		int				false	"Language id"
// @Param			sort_by				query		string			false	"Sorting by a column"
// @Param			limit				query		int				false	"Limit returned number of resources"
// @Param			offset				query		int				false	"Offset returned resources"
// @Param			extend				query		bool			false	"Return extended book information"
// @Param			author.id			query		int				false	"If extend=true - Author id"
// @Param			author.first_name	query		string			false	"If extend=true - Author first name"
// @Param			author.last_name	query		string			false	"If extend=true - Author last name"
// @Success		200					{array}		models.Book		"OK - Fetched books"
// @Failure		400					{object}	models.Error	"Bad Request - Invalid input"
// @Failure		401					{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		500					{object}	models.Error	"Internal Server Error"
// @Router			/books [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetBooks(c *gin.Context) {
	params := c.Request.URL.Query()
	extend := c.DefaultQuery("extend", "false")

	var (
		books any
		err   error
	)

	if extend == "true" {
		books, err = h.DB.GetBooksExt(params)
	} else {
		books, err = h.DB.GetBooks(params)
	}

	if errors.Is(err, db.ErrParam) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
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
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		404	{object}	models.Error	"Not Found - No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/books/{id} [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	book, err := h.DB.GetBook(int64(id))
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary		Create a new book
// @Description	Accepts a JSON body to create a new book. Responds with the created book and set `Location` header or an error message.
// @Tags			Books
// @Accept			json
// @Produce		json
// @Param			book	body		models.Book		true	"New Book"
// @Success		201		{object}	models.Book		"Created - Added new book"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Header			201		{string}	Location		"Path of the newly created book"
// @Router			/books [post]
// @Security		ApiKeyAuth
func (h *Handlers) PostBook(c *gin.Context) {
	var newBook models.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newBook.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	id, err := h.DB.InsertBook(newBook)
	if errors.Is(err, db.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
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
// @Success		204		"No content - Updated the book"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/books/{id} [put]
// @Security		ApiKeyAuth
func (h *Handlers) PutBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var newBook models.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newBook.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	if err := h.DB.UpdateWholeBook(int64(id), newBook); err != nil {
		handleDBError(c, err)
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
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/books/{id} [patch]
// @Security		ApiKeyAuth
func (h *Handlers) PatchBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var patchBook models.Book

	if err := c.BindJSON(&patchBook); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if err := h.DB.UpdateBook(int64(id), patchBook); err != nil {
		handleDBError(c, err)
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
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404	{object}	models.Error	"Not Found -  No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/books/{id} [delete]
// @Security		ApiKeyAuth
func (h *Handlers) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	if err := h.DB.DelBook(int64(id)); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Return allowed operations for books
// @Description	Responds with an empty response body.
// @Tags			Books
// @Success		204	"No Content - Successfully responded with available options"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Header			204	{string}	Allow			"Allowed operations for the resource"
// @Router			/books [options]
// @Security		ApiKeyAuth
func (h *Handlers) OptionsBooks(c *gin.Context) {
	c.Header("Allow", "GET, POST, OPTIONS")
	c.Status(http.StatusNoContent)
}

// @Summary		Return allowed operations for books
// @Description	Responds with an empty response body.
// @Tags			Books
// @Param			id	path	string	true	"Book id"
// @Success		204	"No Content - Successfully responded with available options"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Header			204	{string}	Allow			"Allowed operations for the resource"
// @Router			/books/{id} [options]
// @Security		ApiKeyAuth
func (h *Handlers) OptionsBook(c *gin.Context) {
	c.Header("Allow", "GET, PUT, PATCH, DELETE, OPTIONS")
	c.Status(http.StatusNoContent)
}
