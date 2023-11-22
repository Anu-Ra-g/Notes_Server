package controllers

import (
	"fmt"
	"net/http"
	"os"

	"notes_application/initializers"
	"notes_application/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// SignUp handles the user registration process
func SignUp(c *gin.Context) {
	var body struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// Parse the request body into the 'body' struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the request body using the validator package
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		fmt.Println("Invalid JSON body")
	}

	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		fmt.Println("Failed to hash Password")
	}

	// Create a new user object
	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}

	// Save the user object to the database
	result := initializers.DB.Create(&user)

	// Check for errors during database operation
	if result.Error != nil {
		fmt.Println("Failed to create user")
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "user created",
	})
}

// Login handles the user login process and issues a JWT token upon successful login
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// Parse the request body into the 'body' struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the request body using the validator package
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		fmt.Println("Invalid JSON body")
	}

	// Retrieve the user from the database based on the provided email
	var user models.User
	fetch := initializers.DB.First(&user, "email = ?", body.Email)

	// Check if the user was not found in the database
	if fetch.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare the hashed password from the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Password",
		})
		return
	}

	// Create a new JWT token with the user ID as a claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
	})

	// Sign the token with the application's secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Respond with the generated JWT token
	c.JSON(http.StatusOK, gin.H{
		"sid": tokenString,
	})
}
