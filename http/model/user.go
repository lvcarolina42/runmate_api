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
}

func (c *CreateUserInput) ToEntity() *entity.User {
	return &entity.User{
		Username:  c.Username,
		Name:      c.Name,
		Email:     c.Email,
		Password:  c.Password,
		Birthdate: c.Birthdate,
		Role:      0,
	}
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
