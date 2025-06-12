package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(isAdmin bool) (string, error) {
	timeNow := time.Now().Unix()
	halfHour := int64(time.Hour/time.Second) >> 1

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":   "server",
			"sub":   "user",
			"exp":   timeNow + halfHour,
			"iat":   timeNow,
			"admin": isAdmin,
		},
	)

	s, err := t.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return s, nil
}

// @Summary		Get a JWT token
// @Description	Return a valid JWT token used for authentication and authorization. Optional boolean admin parameter provides creation of admin access token.
// @Tags			Auth
// @Param			admin	query		bool				false	"Return an admin token"
// @Success		200		{object}	map[string]string	"OK - Response body contains JWT token"
// @Failure		400		{object}	map[string]string	"Bad Request - Invalid parameter value"
// @Failure		500		{object}	map[string]string	"Internal Server Error - Failed to create JWT token"
// @Router			/login [get]
func ReturnToken(c *gin.Context) {
	wantAdmin := c.DefaultQuery("admin", "false")

	if wantAdmin != "false" && wantAdmin != "true" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Podano błędną wartość parametru"})
		return
	}

	boolAdmin := false
	if wantAdmin == "true" {
		boolAdmin = true
	}

	token, err := createToken(boolAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy tworzeniu tokenu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
