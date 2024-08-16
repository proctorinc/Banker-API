CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255),
    email VARCHAR(255),
    passwordHash VARCHAR(255)
);

CREATE TABLE merchants (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    customerId VARCHAR(36) REFERENCES users (id) NOT NULL
);

CREATE TABLE institutions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ownerId VARCHAR(36) REFERENCES users (id) NOT NULL
);

CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    institutionId VARCHAR(36) REFERENCES institutions (id) NOT NULL
);

CREATE TABLE transactions (
    id VARCHAR(36) PRIMARY KEY,
    amount DECIMAL(12, 2) DEFAULT 0,
    ownerId VARCHAR(36) REFERENCES users (id) NOT NULL,
    accountId VARCHAR(36) REFERENCES accounts (id) NOT NULL,
    merchantId VARCHAR(36) REFERENCES merchants (id) NOT NULL
);
