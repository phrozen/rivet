package main

import (
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	//"github.com/tylerb/graceful"
)

func main() {

	app := NewApp()
	app.Initialize()
	defer app.dbs.Close()
	defer app.dba.Close()

	e := echo.New()
	e.Use(mw.Logger())

	e.Get("/echo/*", func(c *echo.Context) error {
		return c.String(200, c.Param("_*"))
	})

	e.Get("/login", app.Login)
	e.Get("/logout", app.Logout)

	s := e.Group("/store")
	s.Use(Auth(app))
	s.Get("/:user", app.List)
	s.Get("/:user/*", app.Read)
	s.Post("/:user/*", app.Create)
	s.Put("/:user/*", app.Update)
	s.Delete("/:user/*", app.Delete)
	s.WebSocket("/:user/websocket", app.WebSocket)

	fmt.Println("Listening on port:", app.port)
	e.Run(":" + app.port)
	//graceful.ListenAndServe(e.Server(":"+app.port), 0)
}
