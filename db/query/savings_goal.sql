-- name: CreateSavingsGoal :one
INSERT INTO savings_goals (
  owner,
  title,
  target_amount,
  current_amount,
  target_date,
  linked_account_id
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSavingsGoal :one
SELECT * FROM savings_goals
WHERE id = $1 LIMIT 1;

-- name: ListSavingsGoals :many
SELECT * FROM savings_goals
WHERE owner = $1
ORDER BY created_at DESC;

-- name: AddToSavingsGoal :one
UPDATE savings_goals
SET current_amount = current_amount + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteSavingsGoal :exec
DELETE FROM savings_goals
WHERE id = $1;
