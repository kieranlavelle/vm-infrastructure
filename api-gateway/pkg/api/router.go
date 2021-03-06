package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	auth "github.com/kieranlavelle/api_gateway/pkg/auth"
)

// RouteRequests takes the application from the path and routes the request
func RouteRequests(ctx *gin.Context, conn *pgxpool.Pool, username string) {

	application := ctx.Param("application")
	var containerName string
	var containerPort int

	rows := getAPIProxy(application, conn)
	err := rows.Scan(&containerName, &containerPort)
	if err != nil {
		log.Printf("error getting address rows: %v\n", err)
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"detail": "no application found",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"detail": "please try again later.",
			})
		}
		return
	}

	proxy := createProxy(ctx, containerName, application, containerPort)
	ctx.Request.Header.Set("X-Authenticated-Userid", username)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

// CreateRoutes creates all of the routes and sets up the database connection
func CreateRoutes() *gin.Engine {

	conn := connectToDatabase()

	router := gin.Default()

	router.Use(corsMiddleware())
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

	router.Any("/:application/*path", func(c *gin.Context) { auth.CheckUser(c, conn, RouteRequests) })

	return router
}
