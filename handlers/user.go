package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"myapp/cache"
	"myapp/database"
	"myapp/models"
	"myapp/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUserHandler(c echo.Context) error {
	var users []models.User
	res := database.DB.Find(&users)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	if res.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No Users Found")
	}

	usersDto := make([]types.UserDto, len(users))
	for index, user := range users {
		usersDto[index] = types.NewUserDto(&user)
	}
	return c.JSON(http.StatusOK, usersDto)
}

func GetUserByIdHandler(c echo.Context) error {
	id := c.Param("id")

	// Check if user is in the cache
	var cachedUser models.User
	err := cache.Get(id, &cachedUser)
	if err == nil {
		return c.JSON(http.StatusOK, types.NewUserDto(&cachedUser))
	}

	// User not in cache, query the database
	var user models.User
	res := database.DB.First(&user, id)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}
	if res.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No User Found")
	}

	// Cache the user data
	if err := cache.Set(id, user, 10*time.Minute); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cache user data")
	}

	return c.JSON(http.StatusOK, types.NewUserDto(&user))
}

func UpdateUserHandler(c echo.Context) error {
	id := c.Param("id")
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	newPassword := c.FormValue("newPassword")
	var user models.User
	res := database.DB.First(&user, id)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}
	if res.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No User Found")
	}

	profilePicturePath, err := UploadProfilePicture(c)
	if err != nil {
		return err
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return echo.ErrUnauthorized
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user.Username = username
	user.Email = email
	user.Password = string(hashedPassword)
	user.ProfilePicture = profilePicturePath

	res = database.DB.Save(&user)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	err = cache.Set(id, user, 10*time.Minute)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Caching error")
	}

	return c.JSON(http.StatusOK, types.NewUserDto(&user))
}
