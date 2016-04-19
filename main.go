package main

import (
	"github.com/labstack/echo"
)

func main() {
	
	app := NewApp(nil)
	defer app.System.DB().Close()
	for _, db := range(app.Store) {
		defer db.DB().Close()
	}

	e := echo.New()

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

	e.Run(":" + app.Config.Port)

}
