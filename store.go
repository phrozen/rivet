package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// #ROUTES

// Lists all the keys in the bucket starting from "offset" key up to "limit" number of keys.
// This returns all the keys joined by "\n" so a simple split will give back an array.
func (app *App) List(c echo.Context) error {
	// Check that we can parse limit, ignore errors
	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = LIMIT
	}
	// Get all keys up to limit
	data := app.Store[c.Param("user")].All("store", c.QueryParam("offset"), int(limit))
	// Join the array with new lines
	return c.String(http.StatusOK, strings.Join(data, "\n"))
}

// Sets a value for key if not exist, returns [201] on creation, [200] if existed,
// [400] if request body can't be read, and [500] if error saving the data.
// Different status codes are useful to know if they key existed previously.
func (app *App) Create(c echo.Context) error {
	// If request body can't be read, return [400]
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Get params (check params?)
	user, key := c.Param("user"), c.Param("_*")
	// Check if key exists already
	exist := app.Store[user].Has("store", key)
	// Save the data
	err = app.Store[user].Set("store", key, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) // [500]
	}
	// [201] if key existed previously
	if !exist {
		return c.String(http.StatusCreated, key)
	}
	return c.String(http.StatusOK, key)
}

// Returns the value for the given key
func (app *App) Read(c echo.Context) error {
	// Look for the value
	data := app.Store[c.Param("user")].Get("store", c.Param("_*"))
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	// Return a string? raw? mimetype?
	return c.String(http.StatusOK, string(data))
}

// Updates the value of an already existing key, (returns [404] if it does not exist?)
func (app *App) Update(c echo.Context) error {
	// If request body can't be read, return [400]
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Get params (check params?)
	user, key := c.Param("user"), c.Param("_*")
	// Return [404] if key does not exist?
	if !app.Store[user].Has("store", key) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	// Saves the data
	err = app.Store[user].Set("store", key, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) // [500]
	}
	return c.String(http.StatusOK, key)
}

// Deletes a key and it's value.
func (app *App) Delete(c echo.Context) error {
	// Get params (check params?)
	user, key := c.Param("user"), c.Param("_*")
	err := app.Store[user].Del("store", key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError) // [500]
	}
	return c.String(http.StatusOK, key)
}

// WIP? Need implementation, operations via websocket should allow all CRUD operations
// with a single persistent connection.
func (app *App) WebSocket(c echo.Context) error {
	return c.String(http.StatusNotFound, "Not implemented")
}
