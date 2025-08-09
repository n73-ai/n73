CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  email VARCHAR(55) NOT NULL UNIQUE,
  active BOOL DEFAULT TRUE,
  role TEXT CHECK (role IN ('user', 'admin')) NOT NULL,
  balance FLOAT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) DEFAULT '',
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    domain VARCHAR(255) DEFAULT '',
    status VARCHAR(255) DEFAULT '',
    port INTEGER DEFAULT 0,
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

CREATE TABLE logs (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  error_scope VARCHAR(255) NOT NULL,
  entity_id VARCHAR(255),
  message TEXT DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
