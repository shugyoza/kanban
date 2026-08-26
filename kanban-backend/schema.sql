-- 1. Boards Table
CREATE TABLE IF NOT EXISTS
  boards (
    id VARCHAR(36) PRIMARY KEY, -- UUID for board identification
    title VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

-- 2. Columns Table (To Do, In Progress, Done)
CREATE TABLE IF NOT EXISTS
  columns (
    id VARCHAR(36) PRIMARY KEY,
    board_id VARCHAR(36) NOT NULL,
    title VARCHAR(50) NOT NULL,
    position INT NOT NULL, -- Tracks the left-to-right order of columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES boards (id) ON DELETE CASCADE
  );

-- 3. Tasks Table (the individual cards within columns)
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) PRIMARY KEY, -- UUID for task identification
    column_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position INT NOT NULL, -- Tracks the top-to-bottom order of tasks within a column
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (column_id) REFERENCES columns (id) ON DELETE CASCADE
  );
