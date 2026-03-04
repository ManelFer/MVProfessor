package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         int       `json:"name"`
	Email        int       `json:"email"`
	PasswordHash int       `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

//dto para cadastro
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

//dto para login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

//resposta token
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
