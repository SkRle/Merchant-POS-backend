package dto

import "time"

// CreateUserReq represents the request body for creating a user
type CreateUserReq struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserReq represents the request body for updating a user
type UpdateUserReq struct {
	Username string `json:"username" validate:"omitempty,min=3,max=32"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

// UserResponse represents the user data returned to the client
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
