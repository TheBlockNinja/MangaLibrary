package dao

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDB() (*sqlx.DB, error) {
	config := mysql.Config{
		AllowNativePasswords: true,
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "mysql:3306",
	}
	db, err := sqlx.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Printf("DB ERROR:%v\n", err.Error())
		return db, err
	}
	err = db.Ping()
	for retries := 0; retries < 20 && err != nil; retries++ {
		fmt.Printf("Attempted to retry connection...%d\n", retries)
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}
	return db, err
}
