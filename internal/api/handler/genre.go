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
