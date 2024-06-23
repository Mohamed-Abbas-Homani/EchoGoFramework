package types

import "github.com/golang-jwt/jwt/v5"

type (
	// SignupPayload ...
	SignupPayload struct {
		Username string `form:"username" binding:"required" validate:"required,min=1,max=32"`
		Email    string `form:"email" binding:"required" validate:"required,email"`
		Password string `form:"password" binding:"required" validate:"required,min=8,max=32"`
	}

	// LoginPayload ...
	LoginPayload struct {
		Email    string `json:"email" form:"email" binding:"required" validate:"required,email"`
		Password string `json:"password" form:"password" binding:"required" validate:"required,min=8,max=32"`
	}

	// JwtCustomClaims are custom claims extending default ones.
	JwtCustomClaims struct {
		UserId uint   `json:"userId"`
		Email  string `json:"email"`
		jwt.RegisteredClaims
	}
)
