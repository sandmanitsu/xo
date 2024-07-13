package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"xo/internal/cache"
	"xo/internal/db"
	"xo/internal/repository"

	"github.com/labstack/echo/v4"
)

func Login(login string, password string, c echo.Context) string {
	if login == "" || password == "" {
		return "Login or Password not entered"
	}

	db := db.InitDbConn()
	defer db.Close()

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
		return "Incorrect login or password"
	}

	cacheUser(user, c)

	repository.IsAuth = true
	return ""
}

// Cached user and set cockie
// expiration time (cookie) = 1 hour
func cacheUser(user repository.User, c echo.Context) {
	token := user.Login + user.HashedPassword
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	cache.Cache[hashedToken] = user

	expiration := time.Now().Add(360 * time.Minute)
	fmt.Println(expiration)
	cookie := http.Cookie{
		Name:    "auth",
		Value:   url.QueryEscape(hashedToken),
		Expires: expiration,
	}

	c.SetCookie(&cookie)
}

// todo Fix 500 when cookie is empty
// todo fix server time
func ReadCookie(c echo.Context) bool {
	cookie, err := c.Cookie("auth")
	if err != nil {
		panic(err)
	}

	if cookie.Value == "" {
		return false
	}

	return true
}
