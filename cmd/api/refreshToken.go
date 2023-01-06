package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type RefreshTokenClaims struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	Token
	claims RefreshTokenClaims
}

func (t *RefreshToken) readToken(r *http.Request) (jwt.Claims, string, error) {
	claims := &t.claims
	cookie, err := r.Cookie("RefreshToken")

	if err != nil {
		return nil, "", err
	}
	return claims, cookie.Value, nil
}
