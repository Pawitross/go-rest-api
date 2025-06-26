package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"pawrest/internal/models"
)

// https://github.com/gin-gonic/gin/issues/814
type AdminBody struct {
	RetAdmin *bool `json:"return_admin_token" binding:"required"`
}

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
// @Description	Return a valid JWT token used for authentication and authorization.
// @Description	Endpoint requires a JSON request body with a `return_admin_token` boolean field. Setting it to `true` returns an admin access token.
// @Tags			Auth
// @Param			admin	body		AdminBody		true	"Return an admin token"
// @Success		200		{object}	models.Token	"OK - Response body contains JWT token"
// @Failure		400		{object}	models.Error	"Bad Request - Invalid parameter value"
// @Failure		500		{object}	models.Error	"Internal Server Error - Failed to create JWT token"
// @Router			/login [post]
func ReturnToken(c *gin.Context) {
	var body AdminBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid JSON in request body"})
		return
	}

	boolAdmin := *body.RetAdmin

	token, err := createToken(boolAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, models.Token{Admin: boolAdmin, Token: token})
}
