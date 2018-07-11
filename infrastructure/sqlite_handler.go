package infrastructure

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satellitex/bbft/config"
)

func NewSQLite3(config *config.Config) *sql.DB {
	conn, err := sql.Open("sqlite3", config.Database.Name)
	if err != nil {
		panic(err)
	}
	return conn
}