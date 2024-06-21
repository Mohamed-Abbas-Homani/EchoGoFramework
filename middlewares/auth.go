package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"myapp/handlers"
	"os"
)

func JWTConfig() echojwt.Config {
	secret := os.Getenv("JWT_SECRET")
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
}
