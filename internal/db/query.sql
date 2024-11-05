-- USERS

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

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
SELECT * FROM accounts AS a
WHERE ownerId = $1
ORDER BY a.name
LIMIT $2 OFFSET @start;

-- name: ListAccountsCount :one
SELECT count(id) FROM accounts AS a
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
WHERE ownerId = $1
ORDER BY date
LIMIT $2 OFFSET @start;

-- name: ListTransactionsCount :one
SELECT count(id) FROM transactions AS a
WHERE ownerId = $1;

-- name: ListTransactionsByAccountIds :many
SELECT t.* FROM transactions AS t, accounts AS a
WHERE t.ownerId = $1 AND agents.id = ANY($2::varchar[])
ORDER BY date;

-- name: ListTransactionsByMerchantId :many
SELECT * FROM transactions
WHERE ownerId = $1 AND merchantId = $2
ORDER BY date;

-- name: ListSpendingTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND amount < 0 AND date BETWEEN @startdate AND @enddate
ORDER BY date;

-- name: ListIncomeTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND amount >= 0 AND date BETWEEN @startdate AND @enddate
ORDER BY date;

-- name: ListAccountSpendingTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount < 0
ORDER BY date;

-- name: ListAccountIncomeTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount >= 0
ORDER BY date;

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
    accountId,
    merchantId
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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

-- MERCHANTS

-- name: GetMerchant :one
SELECT * FROM merchants
WHERE id = $1 and ownerId = $2
LIMIT 1;

-- name: GetMerchantByKey :one
SELECT m.* FROM merchants AS m JOIN merchant_keys AS k ON m.id = k.merchantId
WHERE uploadSource = $1 AND keymatch LIKE $2;

-- name: GetMerchantByName :one
SELECT * FROM merchants
WHERE name = $1;

-- name: GetMerchantBySourceId :one
SELECT * FROM merchants
WHERE sourceId = $1;

-- name: ListMerchants :many
SELECT * FROM merchants
WHERE ownerId = $1
ORDER BY name
LIMIT $2 OFFSET @start;

-- name: ListMerchantsCount :one
SELECT count(id) FROM merchants
WHERE ownerId = $1;

-- name: CreateMerchant :one
INSERT INTO merchants (
    name,
    sourceId,
    ownerId
)
VALUES ($1, $2, $3)
RETURNING *;

-- MERCHANT KEYS

-- name: CreateMerchantKey :one
INSERT INTO merchant_keys (
    keymatch,
    uploadSource,
    merchantId,
    ownerId
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- STATS

-- name: GetTotalSpending :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND amount < 0 AND date BETWEEN @startdate AND @enddate;

-- name: GetTotalIncome :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND amount > 0 AND date BETWEEN @startdate AND @enddate;

-- name: GetAccountSpending :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount < 0;

-- name: GetAccountIncome :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount > 0;
