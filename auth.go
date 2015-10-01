package main

import (
	//"fmt"
	"encoding/hex"
	"github.com/labstack/echo"
	"net/http"
)

// Authentication middleware using Session Tokens
func Auth(app *App) echo.HandlerFunc {

	return func(c *echo.Context) error {
		//Check for session token
		header := c.Request().Header.Get("X-Session-Token")

		// If no token 401, should login
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// Token cannot be decoded or mismatch token length
		// Session Tokens are 32 bytes long (64 hexadecimal characters)
		token, err := hex.DecodeString(header)
		if err != nil || len(token) != 32 {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		// If Session Token is not found (maybe no longer valid due to logout)
		// client gets 401 should login again
		query := c.Param("user")
		session := GetValue(app.dba, []byte("session"), []byte(query))
		if session == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// Finally if tokens don't match
		if string(token) != string(session) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return nil
	}

}

// Creates a new Session Token for the user
func (app *App) NewSession(user []byte) string {
	// Check if session already exists and return the existing token
	token := GetValue(app.dba, []byte("session"), user)
	if token != nil {
		return hex.EncodeToString(token)
	}

	// Generate a new Random Session Token
	token = GenerateRandomKey(32)

	// Update user<->token key, values
	err := SetValue(app.dba, []byte("session"), user, token)

	// Bolt error
	if err != nil {
		return "ERROR" //500
	}

	return hex.EncodeToString(token)
}

func (app *App) DestroySession(user []byte) error {
	// If token doesn't exist just return
	token := GetValue(app.dba, []byte("session"), user)
	if token == nil {
		return nil
	}
	// Delete user<->token
	return DeleteValue(app.dba, []byte("session"), user)
}

// Login route
func (app *App) Login(c *echo.Context) error {
	//Check for Basic Authentication Header
	user, pass, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// If credentials don't match
	if !app.Authenticate(user, pass) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	// Return either a new Session Token or an existing one
	return c.String(http.StatusOK, app.NewSession([]byte(user)))
}

func (app *App) Logout(c *echo.Context) error {
	//Check for Basic Authentication Header
	user, pass, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// If credentials don't match
	if !app.Authenticate(user, pass) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	// Destroy session
	err := app.DestroySession([]byte(user))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.String(http.StatusOK, "")
}
