package services_test

import (
	"database/sql"
	"fmt"
	"myapi/services"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var aSer *services.MyAppService

func TestMain(m *testing.M) {
	dbUser := "docker"
	dbPass := "docker"
	dbDatabase := "sampledb"
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1)/%s?parseTime=true", dbUser, dbPass, dbDatabase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	aSer = services.NewMyAppService(db)
	m.Run()
}

func BenchmarkGetArticleService(b *testing.B) {
	articleID := 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := aSer.GetArticleService(articleID)
		if err != nil {
			b.Error(err)
			break
		}
	}
}
