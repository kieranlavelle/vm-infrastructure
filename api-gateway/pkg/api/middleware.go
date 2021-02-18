package api

import "github.com/gin-gonic/gin"

func setIfNotSet(key, value string, c *gin.Context) {
	if c.GetHeader(key) == "" {
		c.Header(key, value)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		setIfNotSet("Access-Control-Allow-Origin", "*", c)
		setIfNotSet("Access-Control-Allow-Headers", "*", c)
		setIfNotSet("Access-Control-Allow-Methods", "*", c)
		setIfNotSet("Access-Control-Expose-Headers", "*", c)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
