package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerToken := c.GetHeader("Authorization")

		if headerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		userToken, ok := strings.CutPrefix(headerToken, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Bearer prefix in Authorization header"})
			return
		}

		token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			errorMsg := ""

			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				errorMsg = "Malformed token"
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				errorMsg = "Invalid token signature"
			case errors.Is(err, jwt.ErrTokenExpired):
				errorMsg = "Token has expired"
			default:
				errorMsg = "Token verification failed"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to parse token claims"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User authentication data not found"})
			return
		}

		claims, ok := userClaims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to parse token claims"})
			return
		}

		if isAdmin := claims["admin"]; isAdmin != true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have sufficient permissions to access this resource"})
			return
		}

		c.Next()
	}
}
