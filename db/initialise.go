package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

func GetDbConnection() *sqlx.DB {
	// Replace with your database credentials
	dsn := "rabby:rabby123@tcp(localhost:3306)/configDb"

	// Open the database connection
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	return db
}
