package api

// User represents a user registered in the databse
type User struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	hashedPassword string
}

// RegisterUser binds the form data that is passed when a user registers
type RegisterUser struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

type token struct {
	AccessToken string `json:"access_token"`
}

type loggedInUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
