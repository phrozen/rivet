package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"

	// Framework
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomKey generates cryptographic secure randome byte slices.
// Receives just the length of the slice, used for tokens.
func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

// SetUserPassword creates a new user and saves the username and BCrypt hashed password.
func (app *App) SetUserPassword(username, password string) error {
	// Hash the password using Bcrypt with default cost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Save the new user
	return app.System.Set("user", username, hash)
}

// Authenticate username/passsword combination.
func (app *App) Authenticate(username, password string) bool {
	//Check if user exists
	hash := app.System.Get("user", username)
	if hash == nil {
		return false
	}
	//Authenticate passsword
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return false
	}
	return true
}

// NewSession xceates a new Session Token for the user, saves the session and returns a session token.
func (app *App) NewSession(user string) string {
	// Check if session already exists and return the existing token
	token := app.System.Get("session", user)
	if token != nil {
		return hex.EncodeToString(token)
	}

	// Generate a new Random Session Token
	token = GenerateRandomKey(32)

	// Update user<->token key, values
	err := app.System.Set("session", user, token)
	if err != nil {
		return "ERROR" // BoltDB error [500]
	}

	return hex.EncodeToString(token)
}

// DestroySession destroys the session token for the user
func (app *App) DestroySession(user string) error {
	// If token doesn't exist just return
	token := app.System.Get("session", user)
	if token == nil {
		return nil
	}
	// Delete user<->token
	return app.System.Del("session", user)
}

// Auth is the authentication middleware using Session Tokens, every request handled is checked
// for the X-Session-Token Header, decodes the session and authenticates the user to
// it's resources in the store. Returns [400] for a bad token, [401] for non valid
// sessions, and proceeds with request otherwise.
func (app *App) Auth(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		//Check for session token
		header := c.Request().Header.Get("X-Session-Token")

		// If no token [401], should login
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// Token cannot be decoded or mismatch token length [400]
		// Session Tokens are 32 bytes long (64 hexadecimal characters)
		token, err := hex.DecodeString(header)
		if err != nil || len(token) != 32 {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		// If Session Token is not found (maybe no longer valid due to logout)
		// client gets [401] should login again
		session := app.System.Get("session", c.Param("user"))
		if session == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// Finally if tokens don't match [401]
		if string(token) != string(session) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return nil
	}

}

// Login route (/login), credentials are given as Basic Authentication, so be sure to do it over SSL or private networking.
// If authentication is successful, it returns [200] and a 32 byte X-Session-Token in a hex encoded string.
// All authenticated requests should include the X-Session-Token header with this value, if the token is no longer
// valid due to logout or expiration, a [401] usually means you have to call this login endpoint again and obtain a new token.
func (app *App) Login(c echo.Context) error {
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
	return c.String(http.StatusOK, app.NewSession(user))
}

// Logout route (/logout), destroys the X-Session-Token for a given user making it no longer valid.
// To logout the route expects the user credentials as Basic Authentication header to destroy the session.
// All subsecuent requests using the previous X-Session-Token will receive [401] and needs to login again.
func (app *App) Logout(c echo.Context) error {
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
	err := app.DestroySession(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.String(http.StatusOK, "")
}
