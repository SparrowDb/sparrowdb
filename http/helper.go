package http

import (
	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/gin-gonic/gin"
)

// BasicMiddleware middleware of gin to handler OPTIONS
// and cors
func BasicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Server", "SparrowDb")

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		c.Next()
	}
}

// AuthMiddleware middleware to check auth token
func AuthMiddleware(onerr func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _, err := auth.ParseClaimFromRequest(c.Request)

		if err != nil || !token.Valid {
			onerr(c)
		}

		c.Next()
	}
}

func hasPermission(c *gin.Context, role int) bool {
	_, u, err := auth.ParseClaimFromRequest(c.Request)
	if err != nil {
		return false
	}
	return auth.CheckUserPermission(u, role)
}
