package services

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser    = os.Getenv("DB_USER")
	dbPass    = os.Getenv("DB_PASS")
	dbDatabse = os.Getenv("DB_NAME")
	dbConn    = fmt.Sprintf("%s:%s@tcp(127.0.0.1)/%s?parseTime=true", dbUser, dbPass, dbDatabse)
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
