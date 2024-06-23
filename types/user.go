package types

import "myapp/models"

type (
	UserDto struct {
		ID             uint   `json:"id"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profilePicture"`
	}

	UpdateUserPayload struct {
		Username    string `form:"username" binding:"required" validate:"required,min=1,max=32"`
		Email       string `form:"email" binding:"required" validate:"required,email"`
		Password    string `form:"password" binding:"required" validate:"required,min=8,max=32"`
		NewPassword string `form:"password" binding:"required" validate:"required,min=8,max=32"`
	}
)

func NewUserDto(user *models.User) UserDto {
	return UserDto{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}
}
