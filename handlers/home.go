package handlers

import (
	"net/http"
    
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Peace From Mash")
}

