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

// @Summary		Get a list of all genres
// @Description	Responds with a list of all genres as JSON. Optional filtering, sorting and pagination is available through parameters.
// @Tags			Genres
// @Produce		json
// @Param			id		query		string			false	"Genre id"
// @Param			name	query		string			false	"Genre name"
// @Param			sort_by	query		string			false	"Sorting by a column"
// @Param			limit	query		int				false	"Limit returned number of resources"
// @Param			offset	query		int				false	"Offset returned resources"
// @Success		200		{array}		models.Genre	"OK - Fetched genres"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/genres [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetGenres(c *gin.Context) {
	params := c.Request.URL.Query()

	genres, err := h.DB.GetGenres(params)
	if errors.Is(err, db.ErrParam) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// @Summary		Get one genre
// @Description	Responds with the queried genre as JSON or an error message.
// @Tags			Genres
// @Produce		json
// @Param			id	path		int				true	"Genre id"
// @Success		200	{object}	models.Genre	"OK - Fetched genre"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid genre id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		404	{object}	models.Error	"Not Found - No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/genres/{id} [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	genre, err := h.DB.GetGenre(int64(id))
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// @Summary		Create a new genre
// @Description	Accepts a JSON body to create a new genre. Responds with the created genre and set `Location` header or an error message.
// @Tags			Genres
// @Accept			json
// @Produce		json
// @Param			genre	body		models.Genre	true	"New Genre"
// @Success		201		{object}	models.Genre	"Created - Added new genre"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Header			201		{string}	Location		"Path of the newly created genre"
// @Router			/genres [post]
// @Security		ApiKeyAuth
func (h *Handlers) PostGenre(c *gin.Context) {
	var newGenre models.Genre

	if err := c.BindJSON(&newGenre); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newGenre.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	id, err := h.DB.InsertGenre(newGenre)
	if errors.Is(err, db.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	newGenre.Id = id

	location := c.FullPath() + "/" + strconv.FormatInt(newGenre.Id, 10)
	c.Header("Location", location)

	c.JSON(http.StatusCreated, newGenre)
}

// @Summary		Update an existing genre
// @Description	Accepts a JSON body to update a genre. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Genres
// @Accept			json
// @Param			id		path	int				true	"Existing Genre id"
// @Param			genre	body	models.Genre	true	"Updated Genre"
// @Success		204		"No content - Updated the genre"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403		{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404		{object}	models.Error	"Not Found -  No resource found"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/genres/{id} [put]
// @Security		ApiKeyAuth
func (h *Handlers) PutGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var newGenre models.Genre

	if err := c.BindJSON(&newGenre); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newGenre.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	if err := h.DB.UpdateWholeGenre(int64(id), newGenre); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Delete an existing genre
// @Description	Responds with a status code. When an error occurs the response body contains an error message.
// @Tags			Genres
// @Param			id	path	int	true	"Genre id"
// @Success		204	"No Content - Successfully deleted the genre"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid genre id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404	{object}	models.Error	"Not Found -  No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/genres/{id} [delete]
// @Security		ApiKeyAuth
func (h *Handlers) DeleteGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	if err := h.DB.DelGenre(int64(id)); err != nil {
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
