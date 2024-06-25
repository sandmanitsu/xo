package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

type Template struct {
	templates *template.Template
}

func main() {
	e = echo.New()
	e.Static("/", "../public")
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
