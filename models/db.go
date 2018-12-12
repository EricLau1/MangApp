package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

const DRIVER = "mysql"
const USER   = "root"
const PASS   = "@root"
const DBNAME = "dbmanga"

func Connect() *sql.DB {

	URL := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", USER, PASS, DBNAME)

	con, erro := sql.Open(DRIVER, URL)

	if erro != nil {

		panic(erro.Error())

	}

	return con
}