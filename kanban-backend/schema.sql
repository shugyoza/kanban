-- 1. Boards Table
CREATE TABLE IF NOT EXISTS
  boards (
    id VARCHAR(36) PRIMARY KEY, -- UUID for board identification
    title VARCHAR(100) NOT NULL,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
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
    is_archived INT DEFAULT 0, -- 0 for active, 1 for archived (soft deletion). Using int instead of boolean for SQLite compatibility.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (column_id) REFERENCES columns (id) ON DELETE CASCADE
  );


-- 1. Users table to capture login credentials
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- Sessions Tracking Table (For stateful session validation)
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY, -- The secure session UUID string
  user_id TEXT NOT NULL, -- Links back to users table
  expires_at TIMESTAMP NOT NULL, -- Session expiration checkpoint
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Index the session ID space for sub-millisecond retrieval lookups
CREATE INDEX IF NOT EXISTS idx_sessions_id ON sessions(id);


-- Invitations Tracking Table (for account registration)
CREATE TABLE IF NOT EXISTS invitations (
  token TEXT PRIMARY KEY,
  email TEXT NOT NULL, -- Optional: restricts registration to a specific email
  created_by TEXT NOT NULL, -- Links to the admin user who generated it
  expires_at DATETIME NOT NULL, -- Enforcement boundary for time expiration
  used_at DATETIME DEFAULT NULL, -- Acts as a single-use flag (null = valid, timestamp = used)
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(created_by) REFERENCES users(id)
);