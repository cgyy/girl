package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var db *sql.DB

func Init() {
	var err error
	dbn := os.Getenv("dbn")

	if len(dbn) == 0 {
		dbn = "root:123456@/todolist?charset=utf8"
	}

	db, err = sql.Open("mysql", dbn)

	checkErr(err)
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetUser() string {
	user := os.Getenv("USER")
	return user
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
