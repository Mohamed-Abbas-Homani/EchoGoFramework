package handlers

import (
	"myapp/database"
	"myapp/models"
	"net/http"

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
