-- USERS

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (username, email, passwordHash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;

-- ACCOUNTS

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 and ownerId = $2
LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE ownerId = $1;

-- name: UpsertAccount :one
INSERT INTO accounts (
    sourceId,
    uploadSource,
    type,
    name,
    routingNumber,
    updated,
    ownerId
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (sourceId) DO UPDATE
SET
    type = $3,
    name = $4,
    routingNumber = $5,
    updated = $6
-- WHERE ownerId = $7 -- HOW DO WE INCLUDE OWNER ID FOR UPDATE
RETURNING *;

-- TRANSACTIONS

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 and ownerId = $2
LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1;

-- name: UpsertTransaction :one
INSERT INTO transactions (
    sourceId,
    uploadSource,
    amount,
    payeeId,
    payee,
    payeeFull,
    isoCurrencyCode,
    date,
    description,
    type,
    checkNumber,
    updated,
    ownerId,
    accountId
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (sourceId) DO UPDATE
SET
    amount = $3,
    payeeId = $4,
    payee = $5,
    payeeFull = $6,
    isoCurrencyCode = $7,
    date = $8,
    description = $9,
    type = $10,
    checkNumber = $11,
    updated = $12
-- WHERE ownerId = $13 -- HOW DO WE INCLUDE OWNER ID FOR UPDATE
RETURNING *;

-- name: UpdateTransaction :one
UPDATE transactions
SET amount = $3
WHERE id = $1 AND ownerId = $2
RETURNING *;

-- name: DeleteTransaction :one
DELETE FROM transactions
WHERE id = $1
RETURNING *;

-- STATS

-- name: GetTotalSpending :one
SELECT SUM(amount) FROM transactions -- HOW DO WE DEFAULT SUM TO 0?
WHERE ownerId = $1 AND amount > 0;

-- name: GetTotalIncome :one
SELECT SUM(amount) FROM transactions -- HOW DO WE DEFAULT SUM TO 0?
WHERE ownerId = $1 AND amount < 0;
