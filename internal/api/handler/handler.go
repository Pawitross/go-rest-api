package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

type Handlers struct {
	DB db.DatabaseInterface
}

func handleDBError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, db.ErrNotFound):
		c.JSON(http.StatusNotFound, models.Error{Error: err.Error()})
	case errors.Is(err, db.ErrForeignKey):
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
	default:
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: "An Internal Server Error occurred"})
	}
}
