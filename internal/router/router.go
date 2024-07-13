package router

import (
	"html/template"
	"io"
	"net/http"
	"path"
	"xo/internal/auth"
	"xo/internal/cache"
	"xo/internal/db"
	"xo/internal/repository"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

type Template struct {
	templates *template.Template
}

func Router(echo *echo.Echo) {
	cache.InitCache()

	e = echo

	e.GET("/", GetPlayground)

	e.GET("/login", LoginPage)
	e.POST("/login", Login)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Get page with playground
func GetPlayground(c echo.Context) error {
	if !repository.IsAuth {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	path := path.Join("../", "public", "html", "playground.html")
	tmpl := &Template{
		templates: template.Must(template.ParseGlob(path)),
	}
	e.Renderer = tmpl

	repository.IsAuth = false
	return c.Render(http.StatusOK, "playground", "")
}

// Get login page
func LoginPage(c echo.Context) error {
	path := path.Join("../", "public", "html", "login.html")
	tmpl := &Template{
		templates: template.Must(template.ParseGlob(path)),
	}
	e.Renderer = tmpl

	var LoginMessage struct {
		Message string `json:"message"`
	}

	return c.Render(http.StatusOK, "login", LoginMessage)
}

// Login process
func Login(c echo.Context) error {
	db := db.InitDbConn()
	defer db.Close()

	var LoginMessage struct {
		Message string `json:"message"`
	}

	LoginMessage.Message = auth.Login(c.FormValue("login"), c.FormValue("password"), c)
	if repository.IsAuth {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	return c.Render(http.StatusOK, "login", LoginMessage)
}
