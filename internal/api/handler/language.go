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

// @Summary		Get a list of all languages
// @Description	Responds with a list of all languages as JSON. Optional filtering, sorting and pagination is available through parameters.
// @Tags			Languages
// @Produce		json
// @Param			id		query		string			false	"Language id"
// @Param			name	query		string			false	"Language name"
// @Param			sort_by	query		string			false	"Sorting by a column"
// @Param			limit	query		int				false	"Limit returned number of resources"
// @Param			offset	query		int				false	"Offset returned resources"
// @Success		200		{array}		models.Language	"OK - Fetched languages"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid input"
// @Failure		401		{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		500		{object}	models.Error	"Internal Server Error"
// @Router			/languages [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetLanguages(c *gin.Context) {
	params := c.Request.URL.Query()

	languages, err := h.DB.GetLanguages(params)
	if errors.Is(err, db.ErrParam) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, languages)
}

// @Summary		Get one language
// @Description	Responds with the queried language as JSON or an error message.
// @Tags			Languages
// @Produce		json
// @Param			id	path		int				true	"Language id"
// @Success		200	{object}	models.Language	"OK - Fetched language"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid language id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		404	{object}	models.Error	"Not Found - No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/languages/{id} [get]
// @Security		ApiKeyAuth
func (h *Handlers) GetLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	language, err := h.DB.GetLanguage(int64(id))
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	c.JSON(http.StatusOK, language)
}

// @Summary		Create a new language
// @Description	Accepts a JSON body to create a new language. Responds with the created language and set `Location` header or an error message.
// @Tags			Languages
// @Accept			json
// @Produce		json
// @Param			language	body		models.Language	true	"New Language"
// @Success		201			{object}	models.Language	"Created - Added new language"
// @Failure		400			{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401			{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403			{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		500			{object}	models.Error	"Internal Server Error"
// @Header			201			{string}	Location		"Path of the newly created language"
// @Router			/languages [post]
// @Security		ApiKeyAuth
func (h *Handlers) PostLanguage(c *gin.Context) {
	var newLanguage models.Language

	if err := c.BindJSON(&newLanguage); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newLanguage.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	id, err := h.DB.InsertLanguage(newLanguage)
	if errors.Is(err, db.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
		return
	}

	newLanguage.Id = id

	location := c.FullPath() + "/" + strconv.FormatInt(newLanguage.Id, 10)
	c.Header("Location", location)

	c.JSON(http.StatusCreated, newLanguage)
}

// @Summary		Update an existing language
// @Description	Accepts a JSON body to update a language. Responds with a status code. When an error occurs the response body contains JSON data with the message.
// @Tags			Languages
// @Accept			json
// @Param			id			path	int				true	"Existing Language id"
// @Param			language	body	models.Language	true	"Updated Language"
// @Success		204			"No content - Updated the language"
// @Failure		400			{object}	models.Error	"Bad Request - Invalid input or JSON"
// @Failure		401			{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403			{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404			{object}	models.Error	"Not Found -  No resource found"
// @Failure		500			{object}	models.Error	"Internal Server Error"
// @Router			/languages/{id} [put]
// @Security		ApiKeyAuth
func (h *Handlers) PutLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	var newLanguage models.Language

	if err := c.BindJSON(&newLanguage); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	if newLanguage.IsNotValid() {
		c.JSON(http.StatusBadRequest, models.Error{Error: "One or more required fields are missing or invalid"})
		return
	}

	if err := h.DB.UpdateWholeLanguage(int64(id), newLanguage); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary		Delete an existing language
// @Description	Responds with a status code. When an error occurs the response body contains an error message.
// @Tags			Languages
// @Param			id	path	int	true	"Language id"
// @Success		204	"No Content - Successfully deleted the language"
// @Failure		400	{object}	models.Error	"Bad Request - Invalid language id"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Failure		404	{object}	models.Error	"Not Found -  No resource found"
// @Failure		500	{object}	models.Error	"Internal Server Error"
// @Router			/languages/{id} [delete]
// @Security		ApiKeyAuth
func (h *Handlers) DeleteLanguage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Provided incorrect identifier"})
		return
	}

	if err := h.DB.DelLanguage(int64(id)); err != nil {
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

// @Summary		Return allowed operations for languages
// @Description	Responds with an empty response body.
// @Tags			Languages
// @Success		204	"No Content - Successfully responded with available options"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Header			204	{string}	Allow			"Allowed operations for the resource"
// @Router			/languages [options]
// @Security		ApiKeyAuth
func (h *Handlers) OptionsLanguages(c *gin.Context) {
	c.Header("Allow", "GET, POST, OPTIONS")
	c.Status(http.StatusNoContent)
}

// @Summary		Return allowed operations for languages
// @Description	Responds with an empty response body.
// @Tags			Languages
// @Param			id	path	string	true	"Language id"
// @Success		204	"No Content - Successfully responded with available options"
// @Failure		401	{object}	models.Error	"Unauthorized - Invalid or missing token"
// @Failure		403	{object}	models.Error	"Forbidden - Insufficient permissions"
// @Header			204	{string}	Allow			"Allowed operations for the resource"
// @Router			/languages/{id} [options]
// @Security		ApiKeyAuth
func (h *Handlers) OptionsLanguage(c *gin.Context) {
	c.Header("Allow", "GET, PUT, DELETE, OPTIONS")
	c.Status(http.StatusNoContent)
}
