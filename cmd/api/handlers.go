package main

import (
	"authentication/data"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type RequestPayload struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

type ResponseDataUser struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

func (app *Config) Registration(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	jsonUtils := JsonUtils{}
	auth := &Auth{&app.Models.User}

	err := jsonUtils.readJSON(w, r, &requestPayload)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	var user data.User
	user.Email = strings.ToLower(requestPayload.Email)
	user.FirstName = requestPayload.FirstName
	user.LastName = requestPayload.LastName
	user.Password = requestPayload.Password

	err = auth.CheckUserExists(user.Email)

	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	idUser, err := app.Models.User.Insert(user)
	user.ID = idUser

	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Register user %d", idUser),
		Data: ResponseDataUser{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	jsonUtils.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	token := Token{}
	jsonUtils := JsonUtils{}
	auth := &Auth{&app.Models.User}

	err := jsonUtils.readJSON(w, r, &requestPayload)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	email := strings.ToLower(requestPayload.Email)
	password := requestPayload.Password

	userData, err := auth.CheckLogin(email, password)

	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	idSession, version, err := app.Models.Session.Insert(email)

	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	accessClaims := &AccessTokenClaims{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
	}

	accessExpirationTime, accessTokenString, err := token.CreateAccessToken(accessClaims)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	refreshClaims := &RefreshTokenClaims{
		ID:      idSession,
		Version: version,
		Email:   email,
	}

	refreshExpirationTime, refreshTokenString, err := token.CreateRefreshToken(refreshClaims, true)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "AccessToken",
		Value:   accessTokenString,
		Expires: accessExpirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshTokenString,
		Expires: refreshExpirationTime,
	})

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintln("Logged"),
	}

	jsonUtils.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) UserData(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	jsonUtils := JsonUtils{}

	err := jsonUtils.readJSON(w, r, &requestPayload)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	ctx := r.Context()

	claims := ctx.Value(ctxTokenClaim{}).(*AccessTokenClaims)

	userData, err := app.Models.User.GetByEmail(claims.Email)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintln("OK"),
		Data: ResponseDataUser{
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
		},
	}

	jsonUtils.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	token := Token{}
	jsonUtils := JsonUtils{}

	err := jsonUtils.readJSON(w, r, &requestPayload)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	ctx := r.Context()

	claims := ctx.Value(ctxTokenClaim{}).(*RefreshTokenClaims)

	session, err := app.Models.Session.GetByID(claims.ID)

	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	if session.Email != claims.Email || session.Version != claims.Version {
		jsonUtils.errorJSON(w, errors.New("invalid refresh token"), http.StatusUnauthorized)
		app.Models.Session.Delete(claims.ID)
		return
	}

	userData, err := app.Models.User.GetByEmail(claims.Email)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	newVersion := app.Models.Session.UpdateVersion(claims.ID, claims.Version)

	accessClaims := &AccessTokenClaims{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
	}

	accessExpirationTime, accessTokenString, err := token.CreateAccessToken(accessClaims)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	refreshClaims := &RefreshTokenClaims{
		ID:               claims.ID,
		Version:          newVersion,
		Email:            claims.Email,
		RegisteredClaims: claims.RegisteredClaims,
	}

	_, refreshTokenString, err := token.CreateRefreshToken(refreshClaims, false)
	if err != nil {
		jsonUtils.errorJSON(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "AccessToken",
		Value:   accessTokenString,
		Expires: accessExpirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshTokenString,
		Expires: claims.RegisteredClaims.ExpiresAt.Time,
	})

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintln("OK"),
	}

	jsonUtils.writeJSON(w, http.StatusAccepted, payload)

}
