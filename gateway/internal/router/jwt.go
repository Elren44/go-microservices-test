package router

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func ParseAccessToken(tokenStr, secret string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token" + err.Error())
	}
	return claims, nil
}

func ParseRefreshToken(tokenStr, secret string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

func GetUserIDFromToken(tokenStr, secret string) (int, error) {
	claims, err := ParseAccessToken(tokenStr, secret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
