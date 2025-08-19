package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{ secret string }

type JWTOption func(*JWTService)

func NewJWT(secret string, _ ...JWTOption) *JWTService { return &JWTService{secret: secret} }

func (j *JWTService) Generate(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.secret))
}

func (j *JWTService) Secret() []byte { return []byte(j.secret) }
