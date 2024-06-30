package router

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io"
	"net/http"
	"path"
	"xo/internal/db"
	"xo/internal/repository"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

type Template struct {
	templates *template.Template
}

func Router(echo *echo.Echo) {
	e = echo

	e.GET("/", GetPlayground)

	e.GET("/login", LoginPage)
	e.POST("/login", Login)
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

func Login(c echo.Context) error {
	db := db.InitDbConn()
	defer db.Close()

	var LoginMessage struct {
		Message string `json:"message"`
	}

	login := c.FormValue("login")
	password := c.FormValue("password")
	if login == "" || password == "" {
		LoginMessage.Message = "Login or Password ot entered"

		return c.Render(http.StatusOK, "login", LoginMessage)
	}

	var user repository.User
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	row := db.QueryRow(
		"SELECT id, login, hashed_password FROM `users` WHERE login = ? and hashed_password = ?",
		login,
		hashedPass,
	)
	err := row.Scan(&user.Id, &user.Login, &user.HashedPassword)
	if err != nil {
		// err = fmt.Errorf("failed to query data: %w", err)
		LoginMessage.Message = "Incorrect login or password"

		return c.Render(http.StatusOK, "login", LoginMessage)
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
