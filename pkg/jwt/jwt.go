package jwt

import (
	"fmt"
	"time"

	"github.com/Hilaladiii/aureus/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtItf interface {
	CreateToken(userID string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type Jwt struct {
	SecretKey  string
	ExpireTime time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewJwt(cfg config.Env) (JwtItf, error) {
	exp, err := time.ParseDuration(cfg.JwtExpire)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format for expiration time: %v", err)
	}

	return &Jwt{
		SecretKey:  cfg.JwtSecret,
		ExpireTime: exp,
	}, nil
}

func (j *Jwt) CreateToken(userID string) (string, error) {
	if j.ExpireTime <= 0 {
		return "", fmt.Errorf("jwt expire time must be greater than 0")
	}

	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token %v", err)
	}

	return signedToken, nil
}

func (j *Jwt) VerifyToken(tokenString string) (string, error) {
	var claims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed parse token %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
