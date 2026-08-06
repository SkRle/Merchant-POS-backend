package tests

import (
	"errors"
	"testing"
	"time"

	"Merchant-POS-backend/model"
	"Merchant-POS-backend/module/user/dto"
	"Merchant-POS-backend/module/user/services"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	users             map[uint]*model.User
	idSeq             uint
	ReturnErr         bool // Toggle to return errors for all operations
	ReturnErrOnUpdate bool // Toggle to return errors only on Update
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*model.User),
		idSeq: 1,
	}
}

func (m *MockUserRepository) Create(user *model.User) error {
	if m.ReturnErr {
		return errors.New("repo error")
	}
	user.ID = m.idSeq
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[m.idSeq] = user
	m.idSeq++
	return nil
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	if m.ReturnErr {
		return nil, errors.New("repo error")
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("record not found")
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	if m.ReturnErr {
		return nil, errors.New("repo error")
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *MockUserRepository) FindAll() ([]*model.User, error) {
	if m.ReturnErr {
		return nil, errors.New("repo error")
	}
	var result []*model.User
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *MockUserRepository) Update(user *model.User) error {
	if m.ReturnErr || m.ReturnErrOnUpdate {
		return errors.New("repo error")
	}
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("record not found")
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(id uint) error {
	if m.ReturnErr {
		return errors.New("repo error")
	}
	if _, ok := m.users[id]; !ok {
		return errors.New("record not found")
	}
	delete(m.users, id)
	return nil
}

// Tests

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	req := &dto.CreateUserReq{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := userService.CreateUser(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", resp.Username)
	}

	// Duplicate email
	_, err = userService.CreateUser(req)
	if err == nil {
		t.Fatalf("expected error for duplicate email, got none")
	}

	// Repo error on create
	req2 := &dto.CreateUserReq{
		Username: "test2",
		Email:    "test2@example.com",
		Password: "password123",
	}
	mockRepo.ReturnErr = true
	_, err = userService.CreateUser(req2)
	if err == nil {
		t.Fatalf("expected error from repo, got none")
	}
	mockRepo.ReturnErr = false

	// Long password bcrypt error (bcrypt has a 72 char limit but let's just trigger a hash error if possible, actually hard to trigger bcrypt error with standard input, so we might not cover bcrypt error easily without mocking it. The coverage will miss one line `if err != nil` after GenerateFromPassword. We can accept that or simulate extremely large pass. Wait, bcrypt handles up to 72 bytes natively. Let's pass a huge password to trigger bcrypt error)
	largePwd := make([]byte, 73)
	for i := range largePwd {
		largePwd[i] = 'a'
	}
	// bcrypt still hashes it, just truncates or errors. Wait, bcrypt returns error if cost is invalid or length is > 72. Actually length > 72 returns error "bcrypt: password length exceeds 72 bytes"
	req3 := &dto.CreateUserReq{
		Username: "test3",
		Email:    "test3@example.com",
		Password: string(largePwd),
	}
	_, err = userService.CreateUser(req3)
	if err == nil {
		t.Errorf("expected bcrypt error for large password")
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	mockRepo.Create(&model.User{Username: "getuser", Email: "get@example.com"})

	_, err := userService.GetUser(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = userService.GetUser(999)
	if err == nil {
		t.Errorf("expected error for non-existent user, got none")
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	mockRepo.Create(&model.User{Username: "user1", Email: "1@example.com"})

	users, err := userService.GetAllUsers()
	if err != nil || len(users) != 1 {
		t.Errorf("Expected 1 user, got error %v", err)
	}

	mockRepo.ReturnErr = true
	_, err = userService.GetAllUsers()
	if err == nil {
		t.Errorf("Expected error from repo, got none")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	mockRepo.Create(&model.User{Username: "olduser", Email: "old@example.com", Password: "oldpassword"})
	mockRepo.Create(&model.User{Username: "other", Email: "other@example.com"})

	req := &dto.UpdateUserReq{Username: "newuser", Email: "new@example.com", Password: "newpassword"}
	_, err := userService.UpdateUser(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Update with duplicate email
	req2 := &dto.UpdateUserReq{Email: "other@example.com"}
	_, err = userService.UpdateUser(1, req2)
	if err == nil {
		t.Errorf("Expected error for duplicate email, got none")
	}

	// User not found
	_, err = userService.UpdateUser(999, req)
	if err == nil {
		t.Errorf("Expected error for non-existent user")
	}

	// Bcrypt error
	largePwd := make([]byte, 73)
	for i := range largePwd {
		largePwd[i] = 'a'
	}
	req3 := &dto.UpdateUserReq{Password: string(largePwd)}
	_, err = userService.UpdateUser(1, req3)
	if err == nil {
		t.Errorf("Expected bcrypt error")
	}

	// Repo error on update
	mockRepo.ReturnErrOnUpdate = true
	_, err = userService.UpdateUser(1, &dto.UpdateUserReq{Username: "fail"})
	if err == nil {
		t.Errorf("Expected repo error")
	}
	mockRepo.ReturnErrOnUpdate = false

	// Empty fields update
	emptyReq := &dto.UpdateUserReq{}
	_, err = userService.UpdateUser(1, emptyReq)
	if err != nil {
		t.Errorf("Expected no error for empty update, got %v", err)
	}

	// Same email update
	sameEmailReq := &dto.UpdateUserReq{Email: "new@example.com"}
	_, err = userService.UpdateUser(1, sameEmailReq)
	if err != nil {
		t.Errorf("Expected no error for same email update, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	mockRepo.Create(&model.User{Username: "deleteuser", Email: "delete@example.com"})

	err := userService.DeleteUser(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = userService.DeleteUser(999)
	if err == nil {
		t.Errorf("Expected error for non-existent user")
	}
}
