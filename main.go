package main

import (
	"github.com/labstack/echo"
)

func main() {

	app := NewApp(nil)
	defer app.System.DB().Close()
	for _, db := range app.Store {
		defer db.DB().Close()
	}

	e := echo.New()

	e.GET("/login", app.Login)
	e.GET("/logout", app.Logout)

	s := e.Group("/store")
	s.Use(app.Auth)
	s.GET("/:user", app.List)
	s.GET("/:user/*", app.Read)
	s.POST("/:user/*", app.Create)
	s.PUT("/:user/*", app.Update)
	s.DELETE("/:user/*", app.Delete)
	//s.WebSocket("/:user/websocket", app.WebSocket)

	e.Start(":" + app.Config.Port)

}
