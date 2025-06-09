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
