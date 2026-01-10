package controllers

import (
	"blog-api-go/models"
	"blog-api-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// NewAuthController returns a new instance of AuthController.
// It takes a pointer to a *gorm.DB as an argument.
// The returned AuthController instance will use the provided DB instance.
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// RegisterInput is a struct that contains the name, email and password fields.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register registers a new user with the provided name, email and password.
// It first checks if the provided JSON is valid, then hashes the provided password
// and creates a new user in the database with the hashed password.
// If the email already exists in the database, it returns a 400 status with a JSON response containing the error "Email already exists".
// If the registration is successful, it returns a 201 status with a JSON response containing the message "registered".
func (a *AuthController) Register(c *gin.Context) {
	// RegisterInput is a struct that contains the name, email and password fields.
	var input RegisterInput

	// ShouldBindJSON reads the request body, maps JSON into a struct, validates it using tags, and returns an error if anything is wrong, allowing the handler to safely stop.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// HashPassword takes a password string and returns a hashed version of it.
	hash, _ := utils.HashPassword(input.Password)

	// User is a struct that contains the name, email and password fields.
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	// Create creates a new user in the database.
	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registered",
	})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login logs in a user with the provided email and password.
// It first checks if the provided JSON is valid, then checks if a user with the provided email exists in the database.
// If the user does not exist, it returns a 401 status with a JSON response containing the error "invalid credentials".
// If the user exists, it then checks if the provided password matches the hashed password in the database.
// If the passwords do not match, it returns a 401 status with a JSON response containing the error "invalid credentials".
// If the login is successful, it returns a 200 status with a JSON response containing a JWT token for the user.
func (a *AuthController) Login(c *gin.Context) {
	// LoginInput is a struct that contains the email and password fields.
	var input LoginInput

	// ShouldBindJSON reads the request body, maps JSON into a struct, validates it using tags, and returns an error if anything is wrong, allowing the handler to safely stop.
	// ? BindJSON is a function in the gin framework that is used to bind JSON data from the request body to a struct.
	// ? ShouldBindJSON reads the request body, maps JSON into a struct, validates it using tags, and returns an error if anything is wrong, allowing the handler to safely stop.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// User is a struct that contains the name, email and password fields.
	var user models.User
	if err := a.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// CheckPassword takes a password string and a hashed version of it, and returns true if they match, false otherwise.
	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// GenerateToken generates a JWT token for a given user ID.
	token, _ := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
