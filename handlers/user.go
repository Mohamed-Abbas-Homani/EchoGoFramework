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
	var users []types.UserDto
	res := database.DB.Model(models.User{}).Find(&users)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	if res.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No Users Found")
	}

	return c.JSON(http.StatusOK, users)
}

func GetUserByIdHandler(c echo.Context) error {
	id := c.Param("id")

	// Check if user is in the cache
	var user types.UserDto
	err := cache.Get(id, &user)
	if err == nil {
		return c.JSON(http.StatusOK, user)
	}

	// User not in cache, query the database
	res := database.DB.Model(models.User{}).First(&user, id)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "No User Found")
	}

	// Cache the user data
	if err := cache.Set(id, user, 10*time.Minute); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cache user data")
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateUserHandler(c echo.Context) error {
	id := c.Param("id")
	var payload types.UpdateUserPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user models.User
	res := database.DB.First(&user, id)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "No User Found")
	}

	profilePicturePath, err := UploadProfilePicture(c)
	if err != nil {
		return err
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user.Username = payload.Username
	user.Email = payload.Email
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

func DeleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	res := database.DB.First(&user, id)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "No User Found")
	}

	res = database.DB.Delete(&user)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	if err := cache.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cache user data")
	}
	return c.JSON(http.StatusOK, types.NewUserDto(&user))
}
