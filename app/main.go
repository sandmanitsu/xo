package main

import (
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

type Template struct {
	templates *template.Template
}

func main() {
	e = echo.New()
	e.Static("/", "../public")

	e.GET("/", GetPlayground)

	e.Logger.Fatal(e.Start(":8080"))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetPlayground(c echo.Context) error {
	path := path.Join("../", "public", "html", "playground.html")
	tmpl := &Template{
		templates: template.Must(template.ParseGlob(path)),
	}
	e.Renderer = tmpl

	return c.Render(http.StatusOK, "playground", "")
}
