# Note

## 20260820
Unit testing passed.
```txt
backend % go test ./internal/usecase/... -v          
=== RUN   TestGetBoardDetails_Success
--- PASS: TestGetBoardDetails_Success (0.00s)
PASS
ok      backend/internal/usecase        (cached)
```

Application is running
```txt
backend % go run cmd/api/main.go           
2026/08/20 16:57:34 Database connection established successfully.
```

## 20260819
While Go naturally uses PascalsCase and databases use snake_case, JSON keys represent the public contract with FrontEnd. Keeping JSON in camelCase ensures a clean separation of concerns and prevents style mismatches across application stack. 

Database (snake_case)  ──>  Go Backend (PascalCase)  ──>  JSON API & Angular (camelCase)
     "created_at"                CreatedAt                     "createdAt"


Ref.: https://upliftorch.com/tools/json-formatter/en/blog/json-best-practices.html


## 20260814
The best starting point is to build the database layer and core data structures first, followed by the Go API layer, and finally the Angular app. Following a backend-first, data-driven approach ensures architecture remains stable.

Building the database layer first (before the API endpoints) ensures we don't accidentally design an API that is impossible or highly inefficient to query. This is known as Data Modeling. Why it protects app architecture:
1. Database constraints dictate business rules. Avoid writing redundant or conflicting validation logic in Go that DB might reject anyway;
2. Avoid N + 1 Query performance trap, when we can optimize query, e.g. using joins, to fetch data efficiently;
3. Clear separation of "Storage Models" vs. "API Models". What is stored in DB should not always match what is being sent to the front end.

### Phases
Step 1: Simple CRUD Monolith ──> Step 2: Complex Frontend State ──> Step 3: Real-Time WebSockets
(Create, read, update, delete)  (Angular Signals + Drag & Drop)     (Go Goroutines & Channels)

We initialized the project with modular monolith to maximize development velocity and eliminate - 
unnecessary network overhead. However, we decoupled the data access layer, business logic, and -
transport layer using Hexagonal Architecture, to ensure that if any domain module experiences - 
unique scaling bottlenecks in the future, it can be seamlessly broken out into a microservice -
without rewriting the core application.

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
