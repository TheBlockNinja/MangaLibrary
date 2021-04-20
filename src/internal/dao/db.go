package dao

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDB(c *cli.Context) (*sqlx.DB, error) {
	config := mysql.Config{
		AllowNativePasswords: true,
		User:                 c.String("mysql-user"),
		Passwd:               c.String("mysql-password"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", c.String("mysql-host"), c.String("mysql-port")),
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
