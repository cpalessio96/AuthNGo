package main

import (
	"authentication/data"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	user *data.User
}

func (auth *Auth) CheckUserExists(email string) error {
	userData, err := auth.user.GetByEmail(email)

	if userData != nil {
		log.Printf("User already exists: %s", userData.Email)
		return errors.New("user already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (auth *Auth) CheckLogin(email string, password string) (*data.User, error) {
	userData, err := auth.user.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return userData, nil
}
