package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/shared-libs/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println(c.Request.URL)

		// getting auth header value
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// extracting token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// validating the token
		claims, err := auth.ValidateToken(tokenParts[1], os.Getenv("APP_SECRET"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// setting username in context
		c.Set("username", claims["username"])
		fmt.Printf("found-claim: %s", claims["username"])
		c.Next()
	}
}
