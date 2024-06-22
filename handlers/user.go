package handlers

import (
	"encoding/json"
	"log"
	"myapp/cache"
	"myapp/database"
	"myapp/models"
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

	return c.JSON(http.StatusOK, users)
}

func GetUserByIdHandler(c echo.Context) error {
	id := c.Param("id")

	// Check if user is in the cache
	cachedUser, err := cache.Get(id)
	if err == nil && cachedUser != "" {
		var user models.User
		err := json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse cached user")
		}
		log.Println("served from cache")
		return c.JSON(http.StatusOK, user)
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
	log.Println("served from database")

	// Cache the user data
	userJSON, err := json.Marshal(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to marshal user")
	}
	if err := cache.Set(id, userJSON, 10*time.Minute); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
