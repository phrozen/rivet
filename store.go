package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strings"
)

// #ROUTES
func (app *App) List(c *echo.Context) error {
	data := make([]string, 0)
	app.dbs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.Param("user")))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			data = append(data, fmt.Sprintf("%s\n", k))
		}

		return nil
	})
	return c.String(http.StatusOK, strings.Join(data, "\n"))
}

func (app *App) Create(c *echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest)
	}
	exist := FindKey(app.dbs, []byte(c.Param("user")), []byte(c.Param("_*")))
	err = SetValue(app.dbs, []byte(c.Param("user")), []byte(c.Param("_*")), data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !exist {
		return c.String(http.StatusCreated, c.Param("_*"))
	}
	return c.String(http.StatusOK, c.Param("_*"))
}

func (app *App) Read(c *echo.Context) error {
	data := GetValue(app.dbs, []byte(c.Param("user")), []byte(c.Param("_*")))
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.String(http.StatusOK, string(data))
}

func (app *App) Update(c *echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest)
	}
	err = SetValue(app.dbs, []byte(c.Param("user")), []byte(c.Param("_*")), data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.String(http.StatusOK, c.Param("_*"))
}

func (app *App) Delete(c *echo.Context) error {
	err := DeleteValue(app.dbs, []byte(c.Param("user")), []byte(c.Param("_*")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.String(http.StatusOK, c.Param("_*"))
}

func (app *App) WebSocket(c *echo.Context) error {
	return c.String(http.StatusOK, "Socket")
}
