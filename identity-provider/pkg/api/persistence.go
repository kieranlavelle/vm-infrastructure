package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func dbConn(c *gin.Context) *pgx.Conn {
	conn, ok := c.MustGet("conn").(*pgx.Conn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": "please try again later",
		})
	}

	return conn
}

func getUser(c *gin.Context, conn *pgx.Conn, username string) (User, error) {

	user := User{}
	query := `
		SELECT
			username, email, hashed_password
		FROM
			users
		WHERE
			username=$1
	`
	err := conn.QueryRow(
		context.Background(), query, username,
	).Scan(&user.Username, &user.Email, &user.hashedPassword)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "invalid username or password",
			})
		default:
			log.Printf("failed to get user: %v. error: %v\n", username, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "please try again later",
			})
		}
		return user, err
	}
	return user, nil
}

func userInDB(c *gin.Context, conn *pgx.Conn, username, email string) (exists bool, err error) {

	query := `
		SELECT
			count(*) > 0
		FROM
			users
		WHERE
			username=$1 OR email=$2
	`
	err = conn.QueryRow(
		context.Background(), query, username, email,
	).Scan(&exists)

	if err != nil {
		log.Printf("failed checking if user existed: user: %v, email: %v, err: %v",
			username, email, err,
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": "please try again later",
		})
		return
	}
	return
}

func createUser(c *gin.Context, conn *pgx.Conn, user RegisterUser) error {
	query := `
		INSERT INTO users
			(username, email, hashed_password)
		VALUES
			($1, $2, $3)
	`
	_, err := conn.Exec(context.Background(), query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("error inserting user: %v. %v\n", user, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": "an unknown error occoured, please try again later",
		})
		return err
	}

	return nil
}
