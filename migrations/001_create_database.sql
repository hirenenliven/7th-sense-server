-- Create the 7th_sense database
-- Run this as the postgres superuser or a user with CREATE DATABASE privileges
-- Usage: psql -U postgres -f migrations/001_create_database.sql

-- Connect to the default postgres database first
\c postgres

-- Drop database if it exists (for clean setup)
-- Comment this out if you want to preserve existing data
DROP DATABASE IF EXISTS 7th_sense;

-- Create the database
CREATE DATABASE 7th_sense
    WITH OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- Grant privileges (adjust username as needed)
-- GRANT ALL PRIVILEGES ON DATABASE 7th_sense TO postgres;

-- Connect to the new database
\c 7th_sense

-- Verify connection
SELECT current_database();
