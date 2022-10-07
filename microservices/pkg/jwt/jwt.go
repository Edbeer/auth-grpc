package jwt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"guthub.com/Edbeer/microservices/internal/core"
)

// Manager
type Manager struct {
	signingKey string
}

// JWT Manager constructor
func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

// JWT Claims struct
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// Generate JWT token
func (m *Manager) GenerateJWTToken(user *core.User) (string, error) {
	claims := &Claims{
		Email: user.Email,
		ID:    user.Uuid.String(),
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

// Parse access token
func (m *Manager) Parse(accessToken string) (*Claims, error) {
	if accessToken == "" {
		log.Fatal("invalid jwt token")
	}

	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
			}
			secret := []byte(m.signingKey)
			return secret, nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}

	return claims, nil
}