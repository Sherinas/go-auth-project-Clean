package usecase

import (
	"errors"
	"testing"

	"github.com/Sherinas/go-auth-project-Clean/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockJWTService mocks the JWTservice interface
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPassword string
		mockRepoSetup func(*MockUserRepository)
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:          "Valid signup",
			inputName:     "Sherina",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Password123!",
			mockRepoSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "sherinascdlm@gmail.com").Return((*domain.User)(nil), gorm.ErrRecordNotFound)
				m.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedUser: &domain.User{
				Name:  "Sherina",
				Email: "sherinascdlm@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:          "Email already exists",
			inputName:     "Sherina",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Password123!",
			mockRepoSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "sherinascdlm@gmail.com").Return(&domain.User{Email: "sherinascdlm@gmail.com"}, nil)
			},
			expectedUser:  nil,
			expectedError: errors.New("email already exists"),
		},
		{
			name:          "Password too short",
			inputName:     "Sherina",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Pw1!",
			mockRepoSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "sherinascdlm@gmail.com").Return((*domain.User)(nil), gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: errors.New("password must be at least 5 characters long"),
		},
		{
			name:          "Invalid password (no special char)",
			inputName:     "Sherina",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Password123",
			mockRepoSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "sherinascdlm@gmail.com").Return((*domain.User)(nil), gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: errors.New("password must contain uppercase, lowercase, number, and special character"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockJWT := new(MockJWTService)
			tt.mockRepoSetup(mockRepo)

			uc := NewUserusecase(mockRepo, mockJWT)
			user, err := uc.SignUp(tt.inputName, tt.inputEmail, tt.inputPassword)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.NotEmpty(t, user.Password) // Password should be hashed
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSignin(t *testing.T) {
	// Hash a sample password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		inputEmail    string
		inputPassword string
		mockSetup     func(*MockUserRepository, *MockJWTService)
		expectedToken string
		expectedError error
	}{
		{
			name:          "Valid signin",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Password123!",
			mockSetup: func(repo *MockUserRepository, jwt *MockJWTService) {
				repo.On("FindByEmail", "sherinascdlm@gmail.com").
					Return(&domain.User{ID: 1, Email: "sherinascdlm@gmail.com", Password: string(hashedPassword)}, nil)
				jwt.On("GenerateToken", uint(1), "sherinascdlm@gmail.com").
					Return("valid-token", nil)
			},
			expectedToken: "valid-token",
			expectedError: nil,
		},
		{
			name:          "Invalid email",
			inputEmail:    "unknown@gmail.com",
			inputPassword: "Password123!",
			mockSetup: func(repo *MockUserRepository, jwt *MockJWTService) {
				repo.On("FindByEmail", "unknown@gmail.com").
					Return((*domain.User)(nil), errors.New("record not found"))
			},
			expectedToken: "",
			expectedError: errors.New("invalid email or password"),
		},
		{
			name:          "Wrong password",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "WrongPass123!",
			mockSetup: func(repo *MockUserRepository, jwt *MockJWTService) {
				repo.On("FindByEmail", "sherinascdlm@gmail.com").
					Return(&domain.User{ID: 1, Email: "sherinascdlm@gmail.com", Password: string(hashedPassword)}, nil)
			},
			expectedToken: "",
			expectedError: errors.New("invalid email or password"),
		},
		{
			name:          "Token generation fails",
			inputEmail:    "sherinascdlm@gmail.com",
			inputPassword: "Password123!",
			mockSetup: func(repo *MockUserRepository, jwt *MockJWTService) {
				repo.On("FindByEmail", "sherinascdlm@gmail.com").
					Return(&domain.User{ID: 1, Email: "sherinascdlm@gmail.com", Password: string(hashedPassword)}, nil)
				jwt.On("GenerateToken", uint(1), "sherinascdlm@gmail.com").
					Return("", errors.New("token error"))
			},
			expectedToken: "",
			expectedError: errors.New("failed to generate token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockJWT := new(MockJWTService)
			tt.mockSetup(mockRepo, mockJWT)

			uc := NewUserusecase(mockRepo, mockJWT)
			token, err := uc.Signin(tt.inputEmail, tt.inputPassword)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}

			mockRepo.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Valid password", "Password123!", true},
		{"No uppercase", "password123!", false},
		{"No lowercase", "PASSWORD123!", false},
		{"No number", "Password!", false},
		{"No special char", "Password123", false},
		{"Too short", "Pw1!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}
