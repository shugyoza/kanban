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
	"kanban-backend/internal/handler"
	"kanban-backend/internal/repository"
	"kanban-backend/internal/usecase"
	"log"
	"net/http"
	"os"
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

	schemaBytes, err := os.ReadFile("./schema.sql")
	if err != nil {
		log.Fatalf("Critical: Failed to read schema.sql file: %v", err)
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Critical: Failed to execute schema.sql migration: %v", err)
	}

	var seederBytes []byte
	seederBytes, err = os.ReadFile("./seeder.sql")
	if err != nil {
		log.Fatalf("Critical: Failed to read seeder.sql file: %v", err)
	}

	// using = instead of := as we are re-using _ and err vars.
	if _, err = db.Exec(string(seederBytes)); err != nil {
		log.Fatalf("Critical: Failed to execute seeder.sql migration: %v", err)
	}

	log.Println("Database tables initialized from schema.sql and seeded from seeder.sql.")

	// 1. Initialize the innermost layer (DB Repository plugin)
	kanbanRepo := repository.NewSQLBoardRepository(db)

	// 2. Inject the repository into the Business Logic layer (UseCase Interactor)
	kanbanUseCase := usecase.NewKanbanInteractor(kanbanRepo)

	// 3. Inject the UseCase into the Outer Delivery Layer (HTTP Handler Plug)
	kanbanHandler := handler.NewKanbanHandler(kanbanUseCase)

	// Map handler's method to a real web URL path endpoint
	http.HandleFunc("GET /api/boards", kanbanHandler.GetBoard)
	http.HandleFunc("PUT /api/tasks/move", kanbanHandler.MoveTask)
	http.HandleFunc("POST /api/tasks", kanbanHandler.CreateTask)
	http.HandleFunc("DELETE /api/tasks", kanbanHandler.DeleteTask)

	serverPort := ":8080"
	// Fire up the native Go local web server
	log.Printf("Kanban backend server is listening at http://localhost%s\n", serverPort)

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Critical: Web server failed to start %v", err)
	}

}
