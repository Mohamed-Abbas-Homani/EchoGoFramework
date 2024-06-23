package handlers

import (
	"errors"
	"fmt"
	"myapp/database"
	"myapp/models"
	"myapp/types"
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
	var payload types.SignupPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if the email already exists in the database
	var existingUser models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&existingUser).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "Email already exists")
	}

	profilePicturePath, err := UploadProfilePicture(c)
	if err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Create user object
	user := models.User{
		Username:       payload.Username,
		Email:          payload.Email,
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

func LoginHandler(c echo.Context) error {
	var payload types.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if the email exists in the database
	var user models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrUnauthorized
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &types.JwtCustomClaims{
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

	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   t,
		Expires: time.Now().Add(24 * time.Hour),
	})

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
