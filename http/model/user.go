package model

import (
	"time"

	"runmate_api/internal/entity"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthdate time.Time `json:"birthdate"`
	Role      int8      `json:"role"`
}

func NewUserFromEntity(user *entity.User) *User {
	return &User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Birthdate: user.Birthdate,
		Role:      user.Role,
	}
}

type CreateUserInput struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Birthdate time.Time `json:"birthdate"`
	Role      int8      `json:"role"`
}
