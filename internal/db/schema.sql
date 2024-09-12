CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;

DROP TYPE IF EXISTS ROLE;
DROP TYPE IF EXISTS ACCOUNT_TYPE;
DROP TYPE IF EXISTS UPLOAD_SOURCE;
DROP TYPE IF EXISTS TRANSACTION_TYPE;

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
    'CHECK'
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
    uploadSource UPLOAD_SOURCE NOT NULL,
    type ACCOUNT_TYPE NOT NULL,
    name VARCHAR(255) NOT NULL,
    routingNumber VARCHAR(255),
    ownerId UUID REFERENCES users (id) NOT NULL
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sourceId VARCHAR(255) NOT NULL UNIQUE,
    uploadSource UPLOAD_SOURCE NOT NULL,
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
    ownerId UUID REFERENCES users (id) NOT NULL,
    accountId UUID REFERENCES accounts (id) NOT NULL
);

alter table transactions
  add constraint check_min_length check (length(sourceid) >= 1);
