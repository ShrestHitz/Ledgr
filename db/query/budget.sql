-- name: CreateBudget :one
INSERT INTO budgets (
  owner,
  category_id,
  monthly_limit,
  month,
  year
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBudget :one
SELECT * FROM budgets
WHERE id = $1 LIMIT 1;

-- name: GetBudgetByCategoryAndMonth :one
SELECT * FROM budgets
WHERE owner = $1
  AND category_id = $2
  AND month = $3
  AND year = $4
LIMIT 1;

-- name: ListBudgets :many
SELECT * FROM budgets
WHERE owner = $1
  AND month = $2
  AND year = $3
ORDER BY id;

-- name: UpdateBudget :one
UPDATE budgets
SET monthly_limit = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1;
