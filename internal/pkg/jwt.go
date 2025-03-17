package pkg

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTservice interface {
	GenerateToken(userID uint, email string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService() JWTservice {
	return &jwtService{
		secretKey: "your-secret-key",
	}
}

func (j *jwtService) GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		log.Println("tttt", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid or expired")
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		log.Println("Claims type: %T, value: %v\n", token.Claims, token.Claims)
		return nil, errors.New("invalid token claims format")
	}
	return token, nil
}
