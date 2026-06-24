package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey         = "secret_key"
	defaultTokenDuration = 24 * time.Hour // one day
)

type JWTClaims struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, email string, name string) (string, error)
	// ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretkey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string) JWTService {
	if secretKey == "" {
		secretKey = jwtSecretKey
	}
	return &jwtService{
		secretkey:     secretKey,
		tokenDuration: defaultTokenDuration,
	}
}

func (js *jwtService) GenerateToken(userId uint, email string, name string) (string, error) {
	// create claims
	claims := JWTClaims{
		UserId: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gotickets",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // create token with claims

	tokenStr, err := token.SignedString([]byte(js.secretkey)) // sign token with secret key
	if err != nil {
		return "", nil
	}
	return tokenStr, nil
}

// func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {}
