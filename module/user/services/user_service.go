package services

import (
	"Merchant-POS-backend/model"
	"Merchant-POS-backend/module/user/dto"
	"Merchant-POS-backend/module/user/ports"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *dto.CreateUserReq) (*dto.UserResponse, error) {
	// Check if user already exists
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return mapToUserResponse(user), nil
}

func (s *userService) GetUser(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return mapToUserResponse(user), nil
}

func (s *userService) GetAllUsers() ([]*dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = mapToUserResponse(user)
	}
	return responses, nil
}

func (s *userService) UpdateUser(id uint, req *dto.UpdateUserReq) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		// Check email uniqueness if changed
		if user.Email != req.Email {
			if _, err := s.repo.FindByEmail(req.Email); err == nil {
				return nil, errors.New("email already in use")
			}
			user.Email = req.Email
		}
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return mapToUserResponse(user), nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}

func mapToUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
