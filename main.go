package main

import (
	"database/sql"
	"fmt"
	"log"
	"myapi/controllers"
	"myapi/routers"
	"myapi/services"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	dbUser := os.Getenv("USERNAME")
	dbPass := os.Getenv("USERPASS")
	dbDatabse := os.Getenv("DATABASE")
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1)/%s?parseTime=true", dbUser, dbPass, dbDatabse)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("fail to connect db")
		return
	}

	ser := services.NewMyAppService(db)
	con := controllers.NewMyAppController(ser)

	r := routers.NewRouter(con)

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
