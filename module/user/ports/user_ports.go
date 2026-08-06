package ports

import (
	"Merchant-POS-backend/model"
	"Merchant-POS-backend/module/user/dto"
)

// UserRepository defines the database operations for User
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

// UserService defines the business logic operations for User
type UserService interface {
	CreateUser(req *dto.CreateUserReq) (*dto.UserResponse, error)
	GetUser(id uint) (*dto.UserResponse, error)
	GetAllUsers() ([]*dto.UserResponse, error)
	UpdateUser(id uint, req *dto.UpdateUserReq) (*dto.UserResponse, error)
	DeleteUser(id uint) error
}
