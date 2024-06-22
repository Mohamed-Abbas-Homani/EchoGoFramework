package handlers

import (
	"errors"
	"fmt"
	"myapp/database"
	"myapp/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUpHandler(c echo.Context) error {
	// Parse form values
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Check if the email already exists in the database
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "Email already exists")
	}

	profilePicturePath, err := UploadProfilePicture(c)
	if err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Create user object
	user := models.User{
		Username:       username,
		Email:          email,
		Password:       string(hashedPassword),
		ProfilePicture: profilePicturePath,
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// Return success response
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
		"userId":  fmt.Sprint(user.ID),
	})
}

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	UserId uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func LoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Check if the email exists in the database
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrUnauthorized
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
