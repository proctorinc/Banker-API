-- USERS

-- name: CreateUser :one
INSERT INTO users (
    username, email, passwordHash
) VALUES (
    $1, $2, $3
)
RETURNING id, username, email;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- TRANSACTIONS

-- name: CreateTransaction :one
INSERT INTO transactions (
    amount, ownerId, accountId, merchantId
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE ownerId = $1 LIMIT 1;

-- name: GetUserTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1;

-- name: GetAccountTransactions :many
SELECT * FROM transactions
WHERE accountId = $1;

-- name: GetMerchantTransactions :many
SELECT * FROM transactions
WHERE merchantId = $1;

-- INSTITUTIONS

-- name: CreateInstitution :one
INSERT INTO institutions (
    name, ownerId
) VALUES (
    $1, $2
)
RETURNING *;

-- ACCOUNTS

-- name: CreateAccount :one
INSERT INTO accounts (
    name, institutionId
) VALUES (
    $1, $2
)
RETURNING *;

-- MERCHANTS

-- name: CreateMerchant :one
INSERT INTO merchants (
    name, customerId
) VALUES (
    $1, $2
)
RETURNING *;
