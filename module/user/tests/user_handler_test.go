package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"Merchant-POS-backend/module/user/dto"
	"Merchant-POS-backend/module/user/handler"

	"github.com/gofiber/fiber/v2"
)

type MockUserService struct {
	CreateUserFunc  func(req *dto.CreateUserReq) (*dto.UserResponse, error)
	GetUserFunc     func(id uint) (*dto.UserResponse, error)
	GetAllUsersFunc func() ([]*dto.UserResponse, error)
	UpdateUserFunc  func(id uint, req *dto.UpdateUserReq) (*dto.UserResponse, error)
	DeleteUserFunc  func(id uint) error
}

func (m *MockUserService) CreateUser(req *dto.CreateUserReq) (*dto.UserResponse, error) {
	return m.CreateUserFunc(req)
}

func (m *MockUserService) GetUser(id uint) (*dto.UserResponse, error) {
	return m.GetUserFunc(id)
}

func (m *MockUserService) GetAllUsers() ([]*dto.UserResponse, error) {
	return m.GetAllUsersFunc()
}

func (m *MockUserService) UpdateUser(id uint, req *dto.UpdateUserReq) (*dto.UserResponse, error) {
	return m.UpdateUserFunc(id, req)
}

func (m *MockUserService) DeleteUser(id uint) error {
	return m.DeleteUserFunc(id)
}

func setupApp(service *MockUserService) *fiber.App {
	app := fiber.New()
	h := handler.NewUserHandler(service)
	app.Post("/users", h.CreateUser)
	app.Get("/users", h.GetAllUsers)
	app.Get("/users/:id", h.GetUser)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
	return app
}

func TestHandler_CreateUser(t *testing.T) {
	mockSvc := &MockUserService{
		CreateUserFunc: func(req *dto.CreateUserReq) (*dto.UserResponse, error) {
			if req.Email == "err@example.com" {
				return nil, errors.New("service error")
			}
			return &dto.UserResponse{ID: 1, Username: req.Username}, nil
		},
	}
	app := setupApp(mockSvc)

	// Success
	body, _ := json.Marshal(dto.CreateUserReq{Username: "test", Email: "test@example.com", Password: "pwd"})
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 201 {
		t.Errorf("Expected 201, got %d", resp.StatusCode)
	}

	// Invalid JSON
	req = httptest.NewRequest("POST", "/users", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

	// Service Error
	body, _ = json.Marshal(dto.CreateUserReq{Username: "test", Email: "err@example.com", Password: "pwd"})
	req = httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_GetUser(t *testing.T) {
	mockSvc := &MockUserService{
		GetUserFunc: func(id uint) (*dto.UserResponse, error) {
			if id == 999 {
				return nil, errors.New("not found")
			}
			return &dto.UserResponse{ID: id}, nil
		},
	}
	app := setupApp(mockSvc)

	// Success
	req := httptest.NewRequest("GET", "/users/1", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Invalid ID
	req = httptest.NewRequest("GET", "/users/abc", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

	// Not Found
	req = httptest.NewRequest("GET", "/users/999", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}

func TestHandler_GetAllUsers(t *testing.T) {
	mockSvc := &MockUserService{
		GetAllUsersFunc: func() ([]*dto.UserResponse, error) {
			return []*dto.UserResponse{{ID: 1}}, nil
		},
	}
	app := setupApp(mockSvc)

	// Success
	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Error
	mockSvc.GetAllUsersFunc = func() ([]*dto.UserResponse, error) {
		return nil, errors.New("db error")
	}
	resp, _ = app.Test(req)
	if resp.StatusCode != 500 {
		t.Errorf("Expected 500, got %d", resp.StatusCode)
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	mockSvc := &MockUserService{
		UpdateUserFunc: func(id uint, r *dto.UpdateUserReq) (*dto.UserResponse, error) {
			if id == 999 {
				return nil, errors.New("user not found")
			}
			if r.Email == "err@example.com" {
				return nil, errors.New("email in use")
			}
			return &dto.UserResponse{ID: id}, nil
		},
	}
	app := setupApp(mockSvc)

	// Success
	body, _ := json.Marshal(dto.UpdateUserReq{Username: "test"})
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Invalid ID
	req = httptest.NewRequest("PUT", "/users/abc", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

	// Invalid body
	req = httptest.NewRequest("PUT", "/users/1", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

	// Not Found
	body, _ = json.Marshal(dto.UpdateUserReq{Username: "test"})
	req = httptest.NewRequest("PUT", "/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}

	// Other Error
	body, _ = json.Marshal(dto.UpdateUserReq{Email: "err@example.com"})
	req = httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	mockSvc := &MockUserService{
		DeleteUserFunc: func(id uint) error {
			if id == 999 {
				return errors.New("user not found")
			}
			return nil
		},
	}
	app := setupApp(mockSvc)

	// Success
	req := httptest.NewRequest("DELETE", "/users/1", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Invalid ID
	req = httptest.NewRequest("DELETE", "/users/abc", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}

	// Not Found
	req = httptest.NewRequest("DELETE", "/users/999", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}
