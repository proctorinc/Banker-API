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

-- name: CountAccounts :one
SELECT count(id) FROM accounts AS a
WHERE ownerId = $1;

-- name: UpsertAccount :one
INSERT INTO accounts (
    sourceId,
    type,
    name,
    routingNumber,
    updated,
    ownerId
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (sourceId) DO UPDATE
SET
    type = $2,
    name = $3,
    routingNumber = $4,
    updated = $5
-- WHERE ownerId = $7 -- HOW DO WE INCLUDE OWNER ID FOR UPDATE
RETURNING *;

-- ACCOUNT SYNC ITEMS

-- name: GetLastSync :one
SELECT * FROM account_sync_items
WHERE accountId = $1
ORDER BY date DESC
LIMIT 1;

-- name: CreateAccountSyncItem :one
INSERT INTO account_sync_items (
    accountId,
    uploadSource
)
VALUES ($1, $2)
RETURNING *;

-- TRANSACTIONS

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 and ownerId = $2
LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1
ORDER BY date DESC
LIMIT $2 OFFSET @start;

-- name: ListTransactionsByDates :many
SELECT * FROM transactions
WHERE ownerId = $1 AND date BETWEEN @startdate AND @enddate
ORDER BY date DESC
LIMIT $2 OFFSET @start;

-- name: ListTransactionsByAccountIds :many
SELECT t.* FROM transactions AS t, accounts AS a
WHERE t.accountId = a.id
    AND a.id::varchar = ANY(@accountIds::varchar[])
ORDER BY date DESC
LIMIT $1 OFFSET @start;

-- name: ListTransactionsByMerchantIds :many
SELECT t.* FROM transactions AS t, merchants AS m
WHERE t.merchantId = m.id
    AND m.id::varchar = ANY(@merchantIds::varchar[])
ORDER BY date DESC
LIMIT $1 OFFSET @start;

-- name: ListSpendingTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1
    AND amount < 0
    AND date BETWEEN @startdate AND @enddate
ORDER BY date DESC
LIMIT $1 OFFSET @start;

-- name: ListIncomeTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1
    AND amount >= 0 AND date BETWEEN @startdate AND @enddate
ORDER BY date DESC
LIMIT $2 OFFSET @start;

-- name: ListAccountSpendingTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount < 0
ORDER BY date DESC
LIMIT $2 OFFSET @start;

-- name: ListAccountIncomeTransactions :many
SELECT * FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount >= 0
ORDER BY date
LIMIT $2 OFFSET @start;

-- name: CountTransactions :one
SELECT count(id) FROM transactions AS a
WHERE ownerId = $1;

-- name: CountTransactionsByDates :one
SELECT count(id) FROM transactions AS a
WHERE ownerId = $1 AND date BETWEEN @startdate AND @enddate;

-- name: CountTransactionsByAccountIds :many
SELECT count(t.id), a.id as accountId FROM transactions AS t, accounts AS a
WHERE t.accountId = a.id
    AND a.id::varchar = ANY(@accountIds::varchar[])
GROUP BY a.id;

-- name: CountTransactionsByMerchantIds :many
SELECT count(t.id), m.id as merchantId FROM transactions AS t, merchants AS m
WHERE t.merchantId = m.id
    AND m.id::varchar = ANY(@merchantIds::varchar[])
GROUP BY m.id;

-- name: CountIncomeTransactions :one
SELECT count(t.id) as merchantId FROM transactions AS t
WHERE ownerId = $1
    AND amount >= 0
    AND date BETWEEN @startdate AND @enddate;

-- name: CountSpendingTransactions :one
SELECT count(t.id) as merchantId FROM transactions AS t
WHERE ownerId = $1
    AND amount < 0
    AND date BETWEEN @startdate AND @enddate;

-- name: UpsertTransaction :one
INSERT INTO transactions (
    sourceId,
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (sourceId) DO UPDATE
SET
    amount = $2,
    payeeId = $3,
    payee = $4,
    payeeFull = $5,
    isoCurrencyCode = $6,
    date = $7,
    description = $8,
    type = $9,
    checkNumber = $10,
    updated = $11
-- WHERE ownerId = $13 -- HOW DO WE INCLUDE OWNER ID FOR UPDATE, NO VALIDATION
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

-- name: ListMerchantsByMerchantIds :many
SELECT m.* FROM transactions AS t, merchants AS m
WHERE t.merchantId = m.id AND m.id::varchar = ANY(@merchantIds::varchar[])
ORDER BY date DESC;

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

-- name: CountMerchants :one
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

-- name: GetNetIncome :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND date BETWEEN @startdate AND @enddate;

-- name: GetAccountSpending :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount < 0;

-- name: GetAccountIncome :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM transactions
WHERE ownerId = $1 AND accountId = $2 AND amount > 0;


-- FUNDS

-- name: CreateFund :one
INSERT INTO funds (type, name, goal, startDate, endDate, ownerId)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListSavingsFunds :many
SELECT * FROM funds
WHERE ownerId = $1 AND type = 'SAVINGS'
ORDER BY name
LIMIT $2 OFFSET @start;

-- name: ListBudgetFunds :many
SELECT * FROM funds
WHERE ownerId = $1 AND type = 'BUDGET'
ORDER BY name
LIMIT $2 OFFSET @start;

-- name: GetFundTotal :one
SELECT COALESCE(SUM(amount), 0) as Sum FROM fund_allocations
WHERE fundId = $1;

-- name: CountSavingsFunds :one
SELECT count(id) FROM funds AS a
WHERE ownerId = $1 AND type = 'SAVINGS';

-- name: CountBudgetFunds :one
SELECT count(id) FROM funds AS a
WHERE ownerId = $1 AND type = 'BUDGET';


-- FUND ALLOCATIONS

-- name: ListFundAllocationsByFundIds :many
SELECT a.* FROM fund_allocations AS a, funds AS f
WHERE a.fundId = f.id
    AND f.id::varchar = ANY(@fundIds::varchar[])
ORDER BY date DESC
LIMIT $1 OFFSET @start;

-- name: CountFundAllocationsByFundId :many
SELECT count(a.id), f.id as fundId FROM fund_allocations AS a, funds AS f
WHERE a.fundId = f.id
    AND f.id::varchar = ANY(@fundIds::varchar[])
GROUP BY f.id;

-- name: GetFundAllocationsStats :one
SELECT
    COALESCE(sum(CASE WHEN a.amount > 0 THEN a.amount ELSE 0 END), 0) as saved,
    COALESCE(sum(CASE WHEN a.amount > 0 THEN a.amount ELSE 0 END), 0) as spent,
    COALESCE(sum(a.amount), 0) as net
FROM fund_allocations AS a, funds AS f
WHERE f.ownerId = $1 AND a.date <= @enddate;


-- MONTHS

-- name: ListMonths :many
SELECT EXTRACT(YEAR FROM t.date) AS year,
       EXTRACT(MONTH FROM t.date) AS month,
       COUNT(*) AS count
FROM transactions AS t
WHERE ownerId = $1
GROUP BY EXTRACT(YEAR FROM t.date), EXTRACT(MONTH FROM t.date)
ORDER BY year DESC, month DESC;
