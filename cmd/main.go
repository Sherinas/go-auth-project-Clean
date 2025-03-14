package main

import (
	"github.com/Sherinas/go-auth-project-Clean/internal/domain"
	"github.com/Sherinas/go-auth-project-Clean/internal/handler"

	"log"

	"github.com/Sherinas/go-auth-project-Clean/internal/pkg"
	"github.com/Sherinas/go-auth-project-Clean/internal/repository"
	"github.com/Sherinas/go-auth-project-Clean/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := "host=localhost user=sherinascdlm password=admin123 dbname=authdb port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	if err := db.AutoMigrate(domain.User{}); err != nil {

		log.Fatal("Faild to migrate")
	}
	jwtService := pkg.NewJWTService()
	userrepo := repository.NewUserRepository(db)
	usecase := usecase.NewUserusecase(userrepo, jwtService)
	authHandler := handler.NewHandler(usecase)

	r := gin.Default()

	r.POST("/signup", authHandler.SignUp)
	r.POST("/signin", authHandler.Signin)
	//	r.GET("/dashboard" auauthHandler.DashBoard)

	r.Run(":8080")
}
