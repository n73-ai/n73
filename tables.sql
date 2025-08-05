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
