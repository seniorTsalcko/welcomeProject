package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	repo      *AuthRepository
	jwtSecret string
}

func NewAuthService(repo *AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) GenerateToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) CheckPassword(user *User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
