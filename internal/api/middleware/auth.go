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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Nie podano tokenu"})
			return
		}

		userToken, ok := strings.CutPrefix(headerToken, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Brak Bearer w nagłówku Authorization"})
			return
		}

		token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			errorMsg := ""

			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				errorMsg = "Zniekształcony token"
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				errorMsg = "Nieprawidłowy podpis tokenu"
			case errors.Is(err, jwt.ErrTokenExpired):
				errorMsg = "Token wygasł"
			default:
				errorMsg = "Błąd weryfikacji"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Błąd z uzyskiwaniem ładunku"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Nie można uzyskać danych użytkownika"})
			return
		}

		claims, ok := userClaims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Błąd z uzyskiwaniem ładunku"})
			return
		}

		if isAdmin := claims["admin"]; isAdmin != true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Nie masz wystarczających uprawnień, aby przeglądać ten zasób"})
			return
		}

		c.Next()
	}
}
