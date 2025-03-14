package handler

import (
	"github.com/Sherinas/go-auth-project-Clean/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.UserUsecase
}

func NewHandler(uc *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (a AuthHandler) SignUp(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required,min=5"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.uc.SignUp(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})

}

func (a *AuthHandler) Signin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=5"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tocken, err := a.uc.Signin(input.Email, input.Password)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signin successfully",
		"tocken":  tocken,
	})

}

func (a *AuthHandler) DashBoard(c *gin.Context) {

	name, exists := c.Get("name")
	email, existsEmail := c.Get("email") // Extracted from middleware

	// If user details are missing, return an error
	if !exists || !existsEmail {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your dashboard!",
		"name":    name,
		"email":   email,
	})
}
