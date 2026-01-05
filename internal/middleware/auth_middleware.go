package middleware

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer([]byte("{}")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
