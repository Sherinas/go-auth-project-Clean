package usecase

import (
	"errors"
	"fmt"

	"github.com/Sherinas/go-auth-project-Clean/internal/domain"
	"github.com/Sherinas/go-auth-project-Clean/internal/pkg"
	"github.com/Sherinas/go-auth-project-Clean/internal/repository"
	"github.com/dlclark/regexp2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase struct {
	repo       repository.UserRepository
	jwtService pkg.JWTservice
}

func NewUserusecase(repo repository.UserRepository, jwt pkg.JWTservice) *UserUsecase {
	return &UserUsecase{repo: repo, jwtService: jwt}
}

func (u *UserUsecase) SignUp(name, email, password string) (*domain.User, error) {
	userExists, err := u.repo.FindByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("DB error:", err)
		return nil, err
	}
	if userExists != nil {
		fmt.Println("Email exists")
		return nil, errors.New("email already exists")
	}
	fmt.Println("Email available, proceeding...")

	if len(password) < 5 {
		fmt.Println("Password too short")
		return nil, errors.New("password must be at least 5 characters long")
	}
	if !isValidPassword(password) {
		fmt.Println("Password validation failed")
		return nil, errors.New("password must contain uppercase, lowercase, number, and special character")
	}
	fmt.Println("Password validated")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Hashing error:", err)
		return nil, err
	}
	fmt.Println("Password hashed")

	user := domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashPassword),
	}
	err = u.repo.Create(&user)
	if err != nil {
		fmt.Println("Create error:", err)
		return nil, err
	}
	fmt.Println("User created")

	return &user, nil
}

func (u *UserUsecase) Signin(email, password string) (string, error) {

	user, err := u.repo.FindByEmail(email)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil

}

func isValidPassword(password string) bool {
	re := regexp2.MustCompile(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`, 0)

	match, _ := re.MatchString(password)
	return match
}
