package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// TokenEntry holds token information
type TokenEntry struct {
	TokenType string
	TokenUUID string
	Username  string
	ExpiresAt time.Time
}

// User holds user information from the database for the user making a request
type User struct {
	Username string
	Email    string
}

// CheckUser checks that the token from the request is valid and passes the user information to the handler
func CheckUser(ctx *gin.Context, conn *pgx.Conn, handler func(*gin.Context, *pgx.Conn, string)) {
	user, err := tokenValid(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Invalid Token.")
		return
	}

	handler(ctx, conn, user)
}

// GetToken generates an access token from the username
func GetToken(ctx *gin.Context) {

	username := ctx.PostForm("username")
	token, err := createJWT(username)
	if err != nil {
		log.Println("Error creating token: ", err)
		ctx.JSON(http.StatusInternalServerError, "Please try again later")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})
	return
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Undexpected signing method")
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TokenValid takes the request context and checks the passed token is valid
func tokenValid(r *http.Request) (string, error) {
	token, err := verifyToken(r)
	if err != nil {
		return "", err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Unable to validate claims")
	}

	var username string = claims["username"].(string)

	return username, nil
}

// CreateJWT creates a new token for the user
func createJWT(username string) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	token := fmt.Sprintf("Bearer %v", accessToken)

	return token, err
}
