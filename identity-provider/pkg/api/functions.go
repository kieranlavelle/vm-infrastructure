package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func connectToDatabase() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func hashPassword(password string) string {
	passwordHash := os.Getenv("PASSWORD_HASH")
	hasher := hmac.New(sha256.New, []byte(passwordHash))
	hasher.Write([]byte(password))

	return hex.EncodeToString(hasher.Sum(nil))
}

func getToken(c *gin.Context, user User) (string, error) {

	tokenData := token{}
	// gatewayAddress := fmt.Sprintf("%v/token", os.Getenv("GATEWAY_ADDRESS"))

	resp, err := http.PostForm("http://api-gateway:8006/token", url.Values{
		"username": {user.Username},
	})
	if err != nil {
		log.Printf("error creating request: %v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		log.Printf("error decoding token: %v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return "", err
	}

	return tokenData.AccessToken, nil
}
