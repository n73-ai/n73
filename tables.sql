DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  email VARCHAR(55) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  active BOOL DEFAULT TRUE,
  token INT NOT NULL,
  verified   BOOL DEFAULT FALSE,
  balance FLOAT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) DEFAULT '',
    status VARCHAR(255) DEFAULT '',
    name TEXT NOT NULL,
    domain VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    project_id VARCHAR(255) NOT NULL,
    role TEXT CHECK (role IN ('user', 'assistant', 'metadata')) NOT NULL,
    content TEXT DEFAULT '',
    model VARCHAR(255) DEFAULT '',
    duration INTEGER DEFAULT 0,
    is_error BOOLEAN DEFAULT false,
    total_cost_usd NUMERIC(10, 6) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
