package tokenserver

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kieranlavelle/api_gateway/pkg/auth"
)

// CreateTokenServer creates the endpoint for giving tokens
func CreateTokenServer() *gin.Engine {

	router := gin.Default()
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

	router.POST("/token", auth.GetToken)

	return router
}
