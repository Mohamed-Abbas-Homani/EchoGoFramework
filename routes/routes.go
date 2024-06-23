package routes

import (
	"github.com/go-playground/validator"
	"myapp/handlers"
	"myapp/middlewares"
	"myapp/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var router *echo.Echo

func InitRoutes() error {
	//Router
	router = echo.New()

	// Use Validator
	router.Validator = &validation.Validator{Validator: validator.New()}

	// Global middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Home Route
	router.GET("/", handlers.HomeHandler)

	// Auth Routes
	authGroup := router.Group("/auth")
	authGroup.POST("/signup", handlers.SignUpHandler)
	authGroup.POST("/login", handlers.LoginHandler)

	// User Routes
	userGroup := router.Group("/user")
	userGroup.Use(echojwt.WithConfig(middlewares.JWTConfig()))
	userGroup.GET("", handlers.GetUserHandler)
	userGroup.GET("/:id", handlers.GetUserByIdHandler)
	userGroup.PUT("/:id", handlers.UpdateUserHandler)
	userGroup.DELETE("/:id", handlers.DeleteUserHandler)
	return nil
}

func Run() {
	router.Logger.Fatal(router.Start(":1323"))
}
