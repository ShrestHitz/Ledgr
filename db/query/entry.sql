-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  category_id,
  amount,
  note
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListEntriesByCategory :many
SELECT * FROM entries
WHERE account_id = $1 AND category_id = $2
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;

-- name: GetAccountSummary :one
SELECT
  COALESCE(SUM(amount) FILTER (WHERE amount > 0), 0)::bigint AS total_income,
  COALESCE(SUM(amount) FILTER (WHERE amount < 0), 0)::bigint AS total_expense,
  COALESCE(SUM(amount), 0)::bigint AS net_balance
FROM entries
WHERE account_id = $1
  AND created_at >= $2
  AND created_at < $3;
