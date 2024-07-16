package router

import (
	"html/template"
	"io"
	"net/http"
	"path"
	"xo/internal/auth"
	"xo/internal/cache"
	"xo/internal/db"

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

	e.POST("/logout", Logout)

	e.GET("/signup", SighupPage)
	e.POST("/signup", Sighup)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Get page with playground
func GetPlayground(c echo.Context) error {
	if !auth.ReadCookie(c) {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	path := path.Join("../", "public", "html", "playground.html")
	tmpl := &Template{
		templates: template.Must(template.ParseGlob(path)),
	}
	e.Renderer = tmpl

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
	if LoginMessage.Message == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	return c.Render(http.StatusOK, "login", LoginMessage)
}

// Logout user
func Logout(c echo.Context) error {
	auth.DeleteCacheAndCookie(c)

	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// Get Sign Up page
func SighupPage(c echo.Context) error {
	path := path.Join("../", "public", "html", "signup.html")
	tmpl := &Template{
		templates: template.Must(template.ParseGlob(path)),
	}
	e.Renderer = tmpl

	return c.Render(http.StatusOK, "signup", "")
}

// Sighup
// POST
func Sighup(c echo.Context) error {
	var LoginMessage struct {
		Message string `json:"message"`
	}

	LoginMessage.Message = auth.AddNewUser(c)
	if LoginMessage.Message == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	return c.Render(http.StatusOK, "signup", LoginMessage)
}
