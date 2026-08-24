// App entry point (wires up DB, UseCases, & Handlers)

/*
Note 1:
The `defer db.Close()` and `defer cancel()` statements ensure that resources are released properly when the main function exits.
This is important for preventing resource leaks (leak db connection or leave orphaned system), especially in long-running applications.
*/
package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite" // SQLite driver
)

func main() {
	// 1. Initialize the pool structure (validates syntax only)

	// SQLite Local File Connection
	db, err := sql.Open("sqlite", "./kanban.db")

	// PostgreSQL Connection. - sslmode=disable is used for local development. In production, you should use sslmode=require or verify-full with proper certificates.
	// db, err := sql.Open("postgres", "user=postgres dbname=kanban sslmode=verify-full")

	if err != nil {
		log.Fatalf("Critical: Invalid database configuration %v", err)
	}
	defer db.Close() // See note 1.

	// 2. Enforce a timeout context to verify the actual connection. 
	defensiveTimingChecks := 3 * time.Second // 3 seconds instead of 2 seconds to allow for some latency (e.g. absorb standard disk or network hiccups without causing a false alarm crash during booting) in the connection.
	ctx, cancel := context.WithTimeout(context.Background(), defensiveTimingChecks)
	defer cancel() // See note 1.

	// 3. Ping the database to ensure the connection is valid and reachable.
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Critical: Unable to reach the database %v", err)
	}

	log.Println("Database connection established successfully.")
}