# README

## Initial Directories Structure
kanban/
│
├── .gitignore                   # The macOS/Go optimized git rules we created
├── README.md                    # Project documentation & architecture log
│
├── todo-backend/                # === GO BACKEND SERVICE ===
│   ├── cmd/
│   │   └── api/
│   │       └── main.go          # App entry point (wires up DB, UseCases, & Handlers)
│   │
│   ├── internal/
│   │   ├── domain/
│   │   │   └── kanban.go        # Business entities, aggregators, & Interface Ports
│   │   │
│   │   ├── repository/
│   │   │   ├── models.go        # Raw DB structural structs (maps to SQL tables)
│   │   │   └── postgres_sqlite.go # SQL data implementation (the LEFT JOIN query)
│   │   │
│   │   ├── usecase/
│   │   │   └── kanban_uc.go     # Business Logic (assembles the nested board tree)
│   │   │
│   │   └── handler/
│   │       └── http.go          # Presentation layer (HTTP requests, REST routers)
│   │
│   ├── schema.sql               # Your SQL Database Schema blueprint
│   ├── go.mod                   # Go module definitions
│   └── go.sum                   # Go checksum tracking
│
└── kanban-frontend/             # === ANGULAR 17+ FRONTEND (Create later) ===
    ├── src/
    │   ├── app/
    │   │   ├── components/      # UI Components (Board, Column, Card)
    │   │   ├── services/        # API Clients / Data fetchers (RxJS/Signals)
    │   │   └── app.config.ts    # Angular global routing and providers
    │   └── main.ts
    ├── angular.json
    └── package.json
