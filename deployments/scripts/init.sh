#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE ${KEYCLOAK_DB_NAME:-keycloak};
    CREATE DATABASE users_db;
    CREATE DATABASE transactions_db;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "users_db" <<-EOSQL
    CREATE TABLE users (
        id UUID PRIMARY KEY,
        username VARCHAR(255),
        email VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "transactions_db" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    CREATE TABLE transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL,
        amount NUMERIC(18, 8) NOT NULL,
        currency VARCHAR(10) NOT NULL,
        transaction_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        categories TEXT[],
        account VARCHAR(100),
        receipt_url TEXT,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
EOSQL