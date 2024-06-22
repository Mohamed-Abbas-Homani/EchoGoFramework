package handlers

import (
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadProfilePicture(c echo.Context) (string, error) {
	profilePicture, err := c.FormFile("profilePicture")
	if err == nil {
		// Open the uploaded file
		src, err := profilePicture.Open()
		if err != nil {
			return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(src)

		// Create destination file
		dst, err := os.Create("uploads/" + profilePicture.Filename)
		if err != nil {
			return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file")
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(dst)

		// Copy file content to destination
		if _, err := io.Copy(dst, src); err != nil {
			return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
		}

		// Set profilePicturePath
		profilePicturePath := "uploads/" + profilePicture.Filename
		return profilePicturePath, nil
	}
	return "", nil
}
