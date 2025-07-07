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

// @Summary		Get a list of all authors
// @Description	Responds with a list of all authors as JSON. Optional filtering, sorting and pagination is available through parameters.
// @Tags			Authors
// @Produce		json
// @Param			id			query		string			false	"Author id"
// @Param			first_name	query		string			false	"Author's first name"
// @Param			last_name	query		string			false	"Author's last name"
// @Param			sort_by		query		string			false	"Sorting by a column"
// @Param			limit		query		int				false	"Limit returned number of resources"
// @Param			offset		query		int				false	"Offset returned resources"
// @Success		200			{array}		models.Author	"OK - Fetched authors"
// @Failure		400			{object}	models.Error	"Bad Request - Invalid input"
// @Failure		401			{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		500			{object}	models.Error	"Internal Server Error"
// @Router			/authors [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetAuthors(c *gin.Context) {
	params := c.Request.URL.Query()

	authors, err := h.DB.GetAuthors(params)
	if errors.Is(err, db.ErrParam) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// @Summary		Get one author
// @Description	Responds with the queried author as JSON or an error message.
// @Tags			Authors
// @Produce		json
// @Param			id	path		int				true	"Author id"
// @Success		200	{object}	models.Author	"OK - Fetched author"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid author id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		404	{object}	models.Error	"Not Found - No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/authors/{id} [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	author, err := h.DB.GetAuthor(int64(id))
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// @Summary		Create a new author
// @Description	Accepts a JSON body to create a new author. Responds with the created author and set `Location` header or an error message.
// @Tags			Authors
// @Accept			json
// @Produce		json
// @Param			author	body		models.Author	true	"New Author"
// @Success		201		{object}	models.Author	"Created - Added new author"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Header			201		{string}	Location		"Path of the newly created author"
// @Router			/authors [post]
// @Security		ApiKeyAuth
func (h *Handlers) PostAuthor(c *gin.Context) {
	var newAuthor models.Author

	if err := c.BindJSON(&newAuthor); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newAuthor.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	id, err := h.DB.InsertAuthor(newAuthor)
	if errors.Is(err, db.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	newAuthor.Id = id

	location := c.FullPath() + "/" + strconv.FormatInt(newAuthor.Id, 10)
	c.Header("Location", location)

	c.JSON(http.StatusCreated, newAuthor)
}

// @Summary		Update an existing author
// @Description	Accepts a JSON body to update a author. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Authors
// @Accept			json
// @Param			id		path	int				true	"Existing Author id"
// @Param			author	body	models.Author	true	"Updated Author"
// @Success		204		"No content - Updated the author"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/authors/{id} [put]
// @Security		ApiKeyAuth
func (h *Handlers) PutAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var newAuthor models.Author

	if err := c.BindJSON(&newAuthor); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newAuthor.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	if err := h.DB.UpdateWholeAuthor(int64(id), newAuthor); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Patch an existing author
// @Description	Accepts a JSON body with patch data to a author. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Authors
// @Accept			json
// @Param			id		path	int				true	"Existing Author id"
// @Param			author	body	models.Author	true	"Patches to the author"
// @Success		204		"No Content - Successfully patched the author"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/authors/{id} [patch]
// @Security		ApiKeyAuth
func (h *Handlers) PatchAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var patchAuthor models.Author

	if err := c.BindJSON(&patchAuthor); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if err := h.DB.UpdateAuthor(int64(id), patchAuthor); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Delete an existing author
// @Description	Responds with a status code. When an error occurs the response body contains an error message.
// @Tags			Authors
// @Param			id	path	int	true	"Author id"
// @Success		204	"No Content - Successfully deleted the author"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid author id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404	{object}	models.Error	"Not Found -  No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/authors/{id} [delete]
// @Security		ApiKeyAuth
func (h *Handlers) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	if err := h.DB.DelAuthor(int64(id)); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.Status(http.StatusNoContent)
}
