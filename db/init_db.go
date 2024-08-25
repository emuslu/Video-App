package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func main() {
	// Open the database connection
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Read and execute the schema SQL file
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		fmt.Println("Error reading schema file:", err)
		return
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		fmt.Println("Error executing schema script:", err)
		return
	}
	fmt.Println("DB initialized")
}
