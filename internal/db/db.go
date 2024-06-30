package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDbConn() *sql.DB {
	db, err := sql.Open("mysql", "admin:admin@tcp(localhost:3306)/my_db")
	if err != nil {
		panic(err)
	}

	return db
}
