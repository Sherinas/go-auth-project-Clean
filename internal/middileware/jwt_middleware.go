package middileware

import (
	"log"
	"net/http"
	"strings"

	"github.com/Sherinas/go-auth-project-Clean/internal/pkg"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(jwtService pkg.JWTservice) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		log.Println(authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		log.Println(parts[0])
		log.Println(len(parts))

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		log.Println("TOCEN,", token)

		claims, ok := token.Claims.(jwt.MapClaims)
		log.Println("ssssss", claims, ok)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims["user_id"])
		c.Set("email", claims["email"])

		c.Next()
	}
}
