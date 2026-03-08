package jwt

import (
	"fmt"
	"time"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtItf interface {
	CreateToken(userID string, role model.Role) (string, error)
	VerifyToken(tokenString string) (*UserClaims, error)
}

type Jwt struct {
	SecretKey  string
	ExpireTime time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string
	Role   model.Role
}

func NewJwt(env config.Env) (JwtItf, error) {
	exp, err := time.ParseDuration(env.JwtExpire)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format for expiration time: %v", err)
	}

	return &Jwt{
		SecretKey:  env.JwtSecret,
		ExpireTime: exp,
	}, nil
}

func (j *Jwt) CreateToken(userID string, role model.Role) (string, error) {
	if j.ExpireTime <= 0 {
		return "", fmt.Errorf("jwt expire time must be greater than 0")
	}

	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token %v", err)
	}

	return signedToken, nil
}

func (j *Jwt) VerifyToken(tokenString string) (*UserClaims, error) {
	var claims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed parse token %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
