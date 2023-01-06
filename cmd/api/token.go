package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ctxTokenClaim struct{}

type IToken interface {
	readToken(*http.Request) (jwt.Claims, string, error)
}

type Token struct {
	IToken
}

func (app *Token) jwtChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonUtils := JsonUtils{}

		claims, tokenString, err := app.readToken(r)

		if err != nil {
			jsonUtils.errorJSON(rw, err, http.StatusBadRequest)
			return
		}

		jwtKey := os.Getenv("SECRET_KEY_JWT")

		isValid, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				jsonUtils.errorJSON(rw, err, http.StatusUnauthorized)
				return
			}
			jsonUtils.errorJSON(rw, err, http.StatusBadRequest)
			return
		}
		if !isValid.Valid {
			jsonUtils.errorJSON(rw, errors.New("jwt not valid"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxTokenClaim{}, claims)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func (app *Token) CreateAccessToken(accessClaims *AccessTokenClaims) (time.Time, string, error) {
	accessExpirationTime := time.Now().Add(5 * time.Minute)

	accessClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	jwtKey := os.Getenv("SECRET_KEY_JWT")

	accessTokenString, err := accessToken.SignedString([]byte(jwtKey))

	return accessExpirationTime, accessTokenString, err
}

func (app *Token) CreateRefreshToken(claims *RefreshTokenClaims, setExpirationDate bool) (time.Time, string, error) {
	refreshExpirationTime := time.Now().Add(24 * time.Hour)

	if setExpirationDate {
		claims.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		}
	}

	jwtKey := os.Getenv("SECRET_KEY_JWT")

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtKey))

	return refreshExpirationTime, refreshTokenString, err

}
