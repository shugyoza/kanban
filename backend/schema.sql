-- 1. Boards Table
CREATE TABLE
  boards (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

-- 2. Columns Table (To Do, In Progress, Done)
CREATE TABLE
  columns (
    id VARCHAR(36) PRIMARY KEY,
    board_id VARCHAR(36) NOT NULL,
    title VARCHAR(50) NOT NULL,
    position INT NOT NULL, -- Tracks the left-to-right order of columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES boards (id) ON DELETE CASCADE
  );