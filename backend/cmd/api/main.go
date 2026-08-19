// App entry point (wires up DB, UseCases, & Handlers)
package main

import (
	"database/sql"
	"log"

	// Swap the driver import when we are ready for Postgres. For now, we are using SQLite for simplicity.
	_ "://sqlite3"
	_ "://postgres"
)

func main() {
	// 1. Initialize the pool structure (validates syntax only)

	// SQLite Local File Connection
	db, err := sql.Open("sqlite3", "./kanban.db")

	// PostgreSQL Connection. - sslmode=disable is used for local development. In production, you should use sslmode=require or verify-full with proper certificates.
	// db, err := sql.Open("postgres", "user=postgres dbname=kanban sslmode=verify-full")

	if err != nil {
		log.Fatalf("Critical: Invalid database configuration %v", err)
	}
	defer db.Close()

	// 2. Enforce a timeout context to verify the actual connection. 3 seconds instead of 2 seconds to allow for some latency in the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 3. Ping the database to ensure the connection is valid and reachable.
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Critical: Unable to reach the database %v", err)
	}

	log.Println("Database connection established successfully.")
}