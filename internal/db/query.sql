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

-- TRANSACTIONS

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1
LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (amount, ownerId)
VALUES ($1, $2)
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
