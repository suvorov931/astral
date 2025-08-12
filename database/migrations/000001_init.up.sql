CREATE SCHEMA IF NOT EXISTS schema_users;

CREATE TABLE IF NOT EXISTS schema_users.users
(
    login TEXT PRIMARY KEY,
    password_hash TEXT
);


