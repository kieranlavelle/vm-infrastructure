package api

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// CreateRoutes connects to the database and creates endpoints for the API.
func CreateRoutes() {

	// form a connection to the database
	connection := connectToDatabase()
	defer connection.Close(context.Background())

	router := gin.Default()

	router.Use(corsMiddleware())
	router.Use(addDatabaseConnection(connection))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Request.Header.Get("X-Authenticated-UserId"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))

	// version the api
	router.POST("/login", loginEndpoint)
	router.POST("/register", registerEndpoint)
	router.GET("/:application/:path", wildcardHandler)

	router.Run("0.0.0.0:8000")
}
