package types

import "myapp/models"

type UserDto struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

func NewUserDto(user *models.User) UserDto {
	return UserDto{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}
}
