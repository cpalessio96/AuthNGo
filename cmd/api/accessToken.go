package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type AccessTokenClaims struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type AccessToken struct {
	Token
	claims AccessTokenClaims
}

func (t *AccessToken) readToken(r *http.Request) (jwt.Claims, string, error) {
	claims := &t.claims
	cookie, err := r.Cookie("AccessToken")

	if err != nil {
		return nil, "", err
	}
	return claims, cookie.Value, nil
}
