package FlatFS

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

func InitSqliteDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func InitMySQLeDB(dataSource string) *sql.DB {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}
