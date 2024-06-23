package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"myapp/types"
	"os"
)

func JWTConfig() echojwt.Config {
	secret := os.Getenv("JWT_SECRET")
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
}
