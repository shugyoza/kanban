-- seed a mock board to test pulling real data
INSERT OR IGNORE INTO boards (id, title) VALUES ('board-kanban-1', 'Kanban Board 1');

-- 2. Columns (Positions 0, 1, 2 build out the horizontal lanes)
INSERT OR IGNORE INTO columns (id, board_id, title, position) VALUES ('col-todo', 'board-kanban-1', 'To Do', 0);
INSERT OR IGNORE INTO columns (id, board_id, title, position) VALUES ('col-progress', 'board-kanban-1', 'In Progress', 1);
INSERT OR IGNORE INTO columns (id, board_id, title, position) VALUES ('col-done', 'board-kanban-1', 'Done', 2);


-- 3. Tasks for Column 1: "To Do"
INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-1', 'col-todo', 'Build Angular Drag & Drop State', 'Implement the optimistic Signal reordering matrix.', 0);

INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-2', 'col-todo', 'Design Go PUT Task Endpoint', 'Create high-performance multi-row index database updater.', 1);

INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-3', 'col-todo', 'Refactor Resume Content', 'Quantify metrics for Sony Music Publishing contributions.', 2);

-- 4. Tasks for Column 2: "In Progress"
INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-4', 'col-progress', 'Master Hexagonal Architecture', 'Learn boundary separation of concerns vs legacy monolith frameworks.', 0);

INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-5', 'col-progress', 'Configure Local Workspace Proxy', 'Bridge Angular port 4200 requests quietly over to Go port 8080.', 1);

-- 5. Tasks for Column 3: "Done"
INSERT OR IGNORE INTO tasks (id, column_id, title, description, position) 
VALUES ('task-6', 'col-done', 'Initialize Git Repository', 'Optimize .gitignore parameters to cleanly ignore macOS spatial logs.', 0);
