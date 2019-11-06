package opendb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// OpenDB opens postgres database
func OpenDB() *sql.DB {
	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		panic(err)
	}
	return db
}
