package utils 

import (
	"regexp"
	"github.com/gin-gonic/gin"
)

func IsSafeString(input string) bool {

	injectionPattern := `(?i)(\bSELECT\b|\bINSERT\b|\bUPDATE\b|\bDELETE\b|\bDROP\b|\bUNION\b|\bEXEC\b|\bALTER\b|\bTRUNCATE\b|\b;|\b--\s)`
	regex := regexp.MustCompile(injectionPattern)
	match := regex.MatchString(input)

	return !match
}

// CORS middleware function
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}