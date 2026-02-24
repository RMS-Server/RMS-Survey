package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/surveyking/server/internal/config"
	"github.com/surveyking/server/internal/dto"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Claims struct {
	User dto.UserInfo `json:"user"`
	jwtlib.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(user dto.UserInfo) (string, error) {
	claims := Claims{
		User: user,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS512, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// ValidateToken parses and validates a JWT string, returning the embedded user info.
func ValidateToken(tokenStr string) (*dto.UserInfo, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &claims.User, nil
}
