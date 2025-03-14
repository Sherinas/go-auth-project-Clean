package repository_test

import (
	"errors"
	"testing"

	"github.com/Sherinas/go-auth-project-Clean/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestFindByEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	expectedUser := &domain.User{Name: "John Doe", Email: "john@example.com"}

	// Mock: If email exists, return user
	mockRepo.On("FindByEmail", "john@example.com").Return(expectedUser, nil)

	user, err := mockRepo.FindByEmail("john@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)

	mockRepo.AssertExpectations(t)
}
func TestFindByEmail_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)

	// Mock: If email doesn't exist, return error
	mockRepo.On("FindByEmail", "notfound@example.com").Return((*domain.User)(nil), errors.New("record not found"))

	user, err := mockRepo.FindByEmail("notfound@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "record not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	newUser := &domain.User{Name: "John Doe", Email: "john@example.com", Password: "Password@123"}

	// Mock: User creation should succeed
	mockRepo.On("Create", newUser).Return(nil)

	err := mockRepo.Create(newUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	newUser := &domain.User{Name: "John Doe", Email: "john@example.com", Password: "Password@123"}

	// Mock: Simulating DB failure
	mockRepo.On("Create", newUser).Return(errors.New("failed to insert user"))

	err := mockRepo.Create(newUser)

	assert.Error(t, err)
	assert.Equal(t, "failed to insert user", err.Error())

	mockRepo.AssertExpectations(t)
}
