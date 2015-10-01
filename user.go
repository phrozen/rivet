package main

import (
	"golang.org/x/crypto/bcrypt"
)

func (app *App) Authenticate(username, password string) bool {
	//Check if user exists
	hash := GetValue(app.dba, []byte("user"), []byte(username))
	if hash == nil {
		return false
	}
	//Authenticate passsword
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return false
	}
	return true
}

func (app *App) SetUserPassword(username, password string) error {
	// Hash the password using Bcrypt with default cost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Save the new user
	return SetValue(app.dba, []byte("user"), []byte(username), hash)
}
