CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS account_sync_items CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS merchants CASCADE;
DROP TABLE IF EXISTS merchant_keys CASCADE;
DROP TABLE IF EXISTS funds CASCADE;
DROP TABLE IF EXISTS fund_allocations CASCADE;

DROP TYPE IF EXISTS ROLE;
DROP TYPE IF EXISTS ACCOUNT_TYPE;
DROP TYPE IF EXISTS UPLOAD_SOURCE;
DROP TYPE IF EXISTS TRANSACTION_TYPE;
DROP TYPE IF EXISTS FUND_TYPE;

CREATE TYPE ROLE AS ENUM (
  'USER',
  'ADMIN'
);

CREATE TYPE ACCOUNT_TYPE AS ENUM (
    'CREDIT',
    'CHECKING',
    'SAVINGS',
    'MONEYMRKT',
    'CREDITLINE',
    'CD'
);

CREATE TYPE UPLOAD_SOURCE AS ENUM (
    'CHASE:CSV_UPLOAD',
    'CHASE:OFX_UPLOAD',
    'PLAID'
);

CREATE TYPE TRANSACTION_TYPE AS ENUM (
    'CREDIT',
    'DEBIT',
    'INT',
    'DIV',
    'FEE',
    'SRVCHG',
    'DEP',
    'ATM',
    'POS',
    'XFER',
    'CHECK',
    'PAYMENT',
    'CASH',
    'DIRECTDEP',
    'DIRECTDEBIT',
    'REPEATPMT',
    'OTHER'
);

CREATE TYPE FUND_TYPE AS ENUM (
    'SAVINGS',
    'BUDGET'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role ROLE DEFAULT 'USER' NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    passwordHash VARCHAR(255) NOT NULL
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sourceId VARCHAR(255) NOT NULL UNIQUE,
    type ACCOUNT_TYPE NOT NULL,
    name VARCHAR(255) NOT NULL,
    routingNumber VARCHAR(255),
    updated DATE NOT NULL DEFAULT NOW(),
    ownerId UUID REFERENCES users (id) NOT NULL
);

CREATE TABLE account_sync_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL DEFAULT NOW(),
    uploadSource UPLOAD_SOURCE NOT NULL,
    accountId UUID REFERENCES accounts (id) NOT NULL
);

CREATE TABLE merchants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    sourceId VARCHAR(255),
    ownerId UUID REFERENCES users (id) NOT NULL
);

CREATE TABLE merchant_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    keymatch VARCHAR(255) NOT NULL,
    uploadSource UPLOAD_SOURCE NOT NULL,
    merchantId UUID REFERENCES merchants (id) NOT NULL,
    ownerId UUID REFERENCES users (id) NOT NULL,
    UNIQUE (ownerId, keymatch)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sourceId VARCHAR(255) NOT NULL UNIQUE,
    amount INT NOT NULL,
    payeeId VARCHAR(255),
    payee VARCHAR(255),
    payeeFull VARCHAR(255),
    isoCurrencyCode VARCHAR(255) DEFAULT 'USD' NOT NULL,
    date DATE NOT NULL,
    description VARCHAR(255) DEFAULT '' NOT NULL,
    type TRANSACTION_TYPE NOT NULL,
    checkNumber VARCHAR(255),
    updated DATE NOT NULL DEFAULT NOW(),
    merchantId UUID REFERENCES merchants (id) NOT NULL,
    ownerId UUID REFERENCES users (id) NOT NULL,
    accountId UUID REFERENCES accounts (id) NOT NULL
);

CREATE TABLE funds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type FUND_TYPE NOT NULL,
    name VARCHAR(255) NOT NULL,
    goal INT NOT NULL DEFAULT 0,
    startDate DATE NOT NULL,
    endDate DATE,
    ownerId UUID REFERENCES users (id) NOT NULL
);

CREATE TABLE fund_allocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    description VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL,
    date DATE NOT NULL,
    ownerId UUID REFERENCES users (id) NOT NULL,
    fundId UUID REFERENCES funds (id) NOT NULL
);
