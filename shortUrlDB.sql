-- PostgreSQL schema for URL Shortener

-- -----------------------
-- Table: users
-- -----------------------
CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- -----------------------
-- Table: short_urls
-- -----------------------
CREATE TABLE public.short_urls (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) REFERENCES public.users(email),
    short_code TEXT,
    actual_url TEXT,
    hits BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- -----------------------
-- Sequences (handled automatically by SERIAL, included for clarity)
-- -----------------------
-- users_id_seq and short_urls_id_seq are automatically created by SERIAL

-- -----------------------
-- Constraints
-- -----------------------
-- users: primary key and unique email are already defined in table
-- short_urls: primary key and foreign key already defined
