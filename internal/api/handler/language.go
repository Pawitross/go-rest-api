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

func (h *Handlers) OptionsLanguages(c *gin.Context) {
	c.Header("Allow", "GET, POST, OPTIONS")
	c.Status(http.StatusNoContent)
}

func (h *Handlers) OptionsLanguage(c *gin.Context) {
	c.Header("Allow", "GET, PUT, DELETE, OPTIONS")
	c.Status(http.StatusNoContent)
}
