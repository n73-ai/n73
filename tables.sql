DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  email VARCHAR(55) NOT NULL UNIQUE,
  active BOOL DEFAULT TRUE,
  role TEXT CHECK (role IN ('user', 'admin')) NOT NULL,
  plan FLOAT DEFAULT 0,
  balance FLOAT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS projects;
CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) DEFAULT '',
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    domain VARCHAR(255) DEFAULT '',
    dev_domain VARCHAR(255) DEFAULT '',
    gh_repo VARCHAR(255) DEFAULT '',
    status VARCHAR(255) DEFAULT '',
    error_msg TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fly_hostname VARCHAR(255) DEFAULT '',

    -- new
    bunny_status TEXT CHECK (bunny_status IN ('storage_zone', 'upload', 'pullzone', 'success')) DEFAULT 'storage_zone' NOT NULL,
    storage_zone_id VARCHAR(255) DEFAULT '',
    storage_zone_region VARCHAR(255) DEFAULT '',
    storage_zone_password VARCHAR(255) DEFAULT '',
    pullzone_id VARCHAR(255) DEFAULT '',

    bunny_eu BOOLEAN DEFAULT false,
    bunny_us BOOLEAN DEFAULT false,
    bunny_asia BOOLEAN DEFAULT false,
    bunny_sa BOOLEAN DEFAULT false,
    bunny_af BOOLEAN DEFAULT false,
    -- end new
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS messages;
CREATE TABLE messages (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    project_id VARCHAR(255) NOT NULL,
    role TEXT CHECK (role IN ('user', 'assistant', 'metadata')) NOT NULL,
    content TEXT DEFAULT '',
    model VARCHAR(255) DEFAULT '',
    duration INTEGER DEFAULT 0,
    is_error BOOLEAN DEFAULT false,
    total_cost_usd NUMERIC(10, 6) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE   
);

DROP TABLE IF EXISTS logs;
CREATE TABLE logs (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  error_scope VARCHAR(255) NOT NULL,
  entity_id VARCHAR(255),
  message TEXT DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
