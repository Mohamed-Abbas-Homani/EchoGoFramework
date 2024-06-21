package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var router *echo.Echo

func InitRoutes() error {
	router = echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.GET("/", handlers.HomeHandler)
	authGroup := router.Group("/auth")
	authGroup.POST("/singup", handlers.SignUpHandler)
	authGroup.POST("/login", handlers.LoginHandler)
	userGroup := router.Group("/user")
	userGroup.Use(echojwt.WithConfig(middlewares.JWTConfig()))
    userGroup.GET("", handlers.GetUserHandler)
	return nil
}

func Run() {
	router.Logger.Fatal(router.Start(":1323"))
}
