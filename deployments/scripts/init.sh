#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE ${KEYCLOAK_DB_NAME:-keycloak};
    CREATE DATABASE users_db;
    \c users_db
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        keycloak_id VARCHAR(255) NOT NULL UNIQUE,
        username VARCHAR(255),
        email VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
EOSQL