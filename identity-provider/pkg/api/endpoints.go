package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginEndpoint(ctx *gin.Context) {

	conn := dbConn(ctx)

	username := ctx.PostForm("username")
	hashedPassword := hashPassword(ctx.PostForm("password"))

	// find the user with that username in the DB
	user, err := getUser(ctx, conn, username)
	if err != nil {
		return
	}

	// Check the supplied password is correct
	if user.hashedPassword == hashedPassword {
		token, err := getToken(ctx, user)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"user":         user,
		})
		return
	}

	ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid username or password"})
}

func registerEndpoint(ctx *gin.Context) {
	conn := dbConn(ctx)
	user := RegisterUser{}

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail": "request body expected fileds [username, email, password]" +
				" as for, data.",
		})
		return
	}

	// check the username and email are not in use
	if exists, err := userInDB(ctx, conn, user.Username, user.Email); err != nil {
		return
	} else if exists {
		ctx.JSON(http.StatusConflict, gin.H{
			"detail": "A user exists with that username or email.",
		})
		return
	}

	if len(user.Username) < 3 || len(user.Username) > 12 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail": "usernames must be between 3 and 12 chars",
		})
		return
	}

	if len(user.Password) < 8 || len(user.Password) > 60 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail": "passwords must be between 8 and 60 chars",
		})
		return
	}

	if !isEmailValid(user.Email) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail": "please provide a valid email message",
		})
		return
	}

	user.Password = hashPassword(user.Password)
	if err := createUser(ctx, conn, user); err != nil {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": loggedInUser{
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

func wildcardHandler(c *gin.Context) {

	log.Printf("Hello from: %v", c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"detail": "wildcards working",
	})
}
