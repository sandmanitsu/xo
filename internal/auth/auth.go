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

	CacheUser(user, c)

	return ""
}

// Cached user and set cockie
// expiration time (cookie) = 6 hour
func CacheUser(user repository.User, c echo.Context) {
	token := user.Login + user.HashedPassword
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	cache.Cache[hashedToken] = user

	expiration := time.Now().Add(360 * time.Minute)
	cookie := http.Cookie{
		Name:    "auth",
		Value:   url.QueryEscape(hashedToken),
		Expires: expiration,
	}

	c.SetCookie(&cookie)
}

// Read a cookie
// todo Fix 500 when cookie is empty
// todo fix server time
func ReadCookie(c echo.Context) bool {
	cookie, err := c.Cookie("auth")
	if err != nil {
		fmt.Println(err)

		return false
	}

	if cookie.Value == "" {
		return false
	}

	_, exist := cache.Cache[cookie.Value]
	if !exist {
		return false
	}

	fmt.Println(cache.Cache)

	return true
}

// Delete cookie and cache
func DeleteCacheAndCookie(c echo.Context) {
	cookie, err := c.Cookie("auth")
	if err != nil {
		fmt.Println(err)
	}

	delete(cache.Cache, cookie.Value)

	emptyCookie := http.Cookie{
		Name:    "auth",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	c.SetCookie(&emptyCookie)
}

// add new user to db
// todo check for same login
func AddNewUser(c echo.Context) string {
	db := db.InitDbConn()
	defer db.Close()

	if c.FormValue("login") == "" || c.FormValue("password") == "" {
		return "Login Or Password not entered"
	}

	hash := md5.Sum([]byte(c.FormValue("password")))
	hashedPass := hex.EncodeToString(hash[:])

	var user repository.User
	user.Login = c.FormValue("login")
	user.HashedPassword = hashedPass

	_, err := db.Exec(`INSERT INTO users (login, hashed_password) values (?, ?)`, user.Login, user.HashedPassword)
	if err != nil {
		fmt.Println(err)

		return fmt.Sprintln("Error: %w", err)
	}

	CacheUser(user, c)

	return ""
}
