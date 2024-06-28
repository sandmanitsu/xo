package main

import (
	"xo/internal/router"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

func main() {

	e = echo.New()
	e.Static("/", "../public")

	router.Router(e)

	e.Logger.Fatal(e.Start(":8080"))
}
