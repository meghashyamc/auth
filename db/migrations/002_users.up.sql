CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4(),
    first_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_digest VARCHAR NOT NULL,
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);