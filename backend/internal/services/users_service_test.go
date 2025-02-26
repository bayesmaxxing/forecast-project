package services_test

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements the repository.UserRepository interface
type MockUserRepository struct {
	users            map[int64]*models.User
	usersByUsername  map[string]*models.User
	nextID           int64
	getUserByIDErr   error
	getUserByNameErr error
	createUserErr    error
	deleteUserErr    error
	validateUserErr  error
	updatePwdErr     error
	listUsersErr     error
}

func NewMockUserRepository() repository.UserRepository {
	return &MockUserRepository{
		users:           make(map[int64]*models.User),
		usersByUsername: make(map[string]*models.User),
		nextID:          1,
	}
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if m.getUserByIDErr != nil {
		return nil, m.getUserByIDErr
	}
	u, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if m.getUserByNameErr != nil {
		return nil, m.getUserByNameErr
	}
	u, exists := m.usersByUsername[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if m.createUserErr != nil {
		return m.createUserErr
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.usersByUsername[user.Username] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	if m.deleteUserErr != nil {
		return m.deleteUserErr
	}
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(m.usersByUsername, user.Username)
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) ValidateUser(ctx context.Context, id int64) (bool, error) {
	if m.validateUserErr != nil {
		return false, m.validateUserErr
	}
	_, exists := m.users[id]
	return exists, nil
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	if m.updatePwdErr != nil {
		return m.updatePwdErr
	}
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	user.Password = password
	return nil
}

func (m *MockUserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	if m.listUsersErr != nil {
		return nil, m.listUsersErr
	}
	var result []*models.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

// Helper function to access the underlying mock
func getMockUserRepo(repo repository.UserRepository) *MockUserRepository {
	return repo.(*MockUserRepository)
}

// Test for GetUserByID
func TestGetUserByID(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user
	testUser := &models.User{
		Username:  "testuser",
		Password:  "password123",
		CreatedAt: time.Now(),
	}
	mockRepo.(*MockUserRepository).CreateUser(ctx, testUser)

	// Test successful retrieval
	user, err := service.GetUserByID(ctx, testUser.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.ID != testUser.ID {
		t.Errorf("Expected user ID %d, got %d", testUser.ID, user.ID)
	}
	if user.Username != testUser.Username {
		t.Errorf("Expected username %s, got %s", testUser.Username, user.Username)
	}

	// Test error case
	getMockUserRepo(mockRepo).getUserByIDErr = errors.New("database error")
	_, err = service.GetUserByID(ctx, testUser.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetUserByUsername
func TestGetUserByUsername(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user
	testUser := &models.User{
		Username:  "testuser",
		Password:  "password123",
		CreatedAt: time.Now(),
	}
	mockRepo.(*MockUserRepository).CreateUser(ctx, testUser)

	// Test successful retrieval
	user, err := service.GetUserByUsername(ctx, testUser.Username)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.ID != testUser.ID {
		t.Errorf("Expected user ID %d, got %d", testUser.ID, user.ID)
	}
	if user.Username != testUser.Username {
		t.Errorf("Expected username %s, got %s", testUser.Username, user.Username)
	}

	// Test error case
	getMockUserRepo(mockRepo).getUserByNameErr = errors.New("database error")
	_, err = service.GetUserByUsername(ctx, testUser.Username)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CreateUser
func TestCreateUser(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Test successful creation
	testUser := &models.User{
		Username: "testuser",
		Password: "password123",
	}
	err := service.CreateUser(ctx, testUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if testUser.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", testUser.ID)
	}

	// Verify password was hashed
	err = bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte("password123"))
	if err != nil {
		t.Error("Password was not properly hashed")
	}

	// Test error case
	getMockUserRepo(mockRepo).createUserErr = errors.New("database error")
	err = service.CreateUser(ctx, testUser)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for VerifyPassword
func TestVerifyPassword(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user with hashed password
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	testUser := &models.User{
		Username:  "testuser",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	getMockUserRepo(mockRepo).CreateUser(ctx, testUser)

	// Test successful verification
	user, err := service.VerifyPassword(ctx, testUser.Username, plainPassword)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.ID != testUser.ID {
		t.Errorf("Expected user ID %d, got %d", testUser.ID, user.ID)
	}

	// Test incorrect password
	_, err = service.VerifyPassword(ctx, testUser.Username, "wrongpassword")
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}

	// Test user not found
	_, err = service.VerifyPassword(ctx, "nonexistentuser", plainPassword)
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
}

// Test for ChangePassword
func TestChangePassword(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user with hashed password
	oldPassword := "oldpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	testUser := &models.User{
		Username:  "testuser",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	getMockUserRepo(mockRepo).CreateUser(ctx, testUser)

	// Test successful password change
	newPassword := "newpassword"
	err := service.ChangePassword(ctx, testUser.ID, oldPassword, newPassword)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify new password was saved and hashed
	userFromRepo, _ := getMockUserRepo(mockRepo).GetUserByID(ctx, testUser.ID)
	err = bcrypt.CompareHashAndPassword([]byte(userFromRepo.Password), []byte(newPassword))
	if err != nil {
		t.Error("New password was not properly saved or hashed")
	}

	// Test incorrect old password
	err = service.ChangePassword(ctx, testUser.ID, "wrongoldpassword", "somepassword")
	if err == nil {
		t.Error("Expected error for wrong old password, got nil")
	}

	// Test user not found
	getMockUserRepo(mockRepo).getUserByIDErr = errors.New("user not found")
	err = service.ChangePassword(ctx, 999, oldPassword, newPassword)
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
}

// Test for DeleteUser
func TestDeleteUser(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user
	testUser := &models.User{
		Username:  "testuser",
		Password:  "password123",
		CreatedAt: time.Now(),
	}
	getMockUserRepo(mockRepo).CreateUser(ctx, testUser)

	// Test successful deletion
	err := service.DeleteUser(ctx, testUser.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was deleted
	_, err = getMockUserRepo(mockRepo).GetUserByID(ctx, testUser.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}

	// Test error case
	getMockUserRepo(mockRepo).deleteUserErr = errors.New("database error")
	err = service.DeleteUser(ctx, testUser.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for ValidateUser
func TestValidateUser(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create a test user
	testUser := &models.User{
		Username:  "testuser",
		Password:  "password123",
		CreatedAt: time.Now(),
	}
	getMockUserRepo(mockRepo).CreateUser(ctx, testUser)

	// Test successful validation
	valid, err := service.ValidateUser(ctx, testUser.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected user to be valid, got invalid")
	}

	// Test invalid user
	valid, err = service.ValidateUser(ctx, 999)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if valid {
		t.Error("Expected user to be invalid, got valid")
	}

	// Test error case
	getMockUserRepo(mockRepo).validateUserErr = errors.New("database error")
	_, err = service.ValidateUser(ctx, testUser.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for ListUsers
func TestListUsers(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)
	ctx := context.Background()

	// Create test users
	user1 := &models.User{Username: "user1", Password: "password1"}
	user2 := &models.User{Username: "user2", Password: "password2"}
	getMockUserRepo(mockRepo).CreateUser(ctx, user1)
	getMockUserRepo(mockRepo).CreateUser(ctx, user2)

	// Test successful listing
	users, err := service.ListUsers(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Test error case
	getMockUserRepo(mockRepo).listUsersErr = errors.New("database error")
	_, err = service.ListUsers(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}