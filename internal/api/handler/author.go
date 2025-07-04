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
