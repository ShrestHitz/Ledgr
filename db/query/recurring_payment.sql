-- name: CreateRecurringPayment :one
INSERT INTO recurring_payments (
  owner,
  title,
  amount,
  category_id,
  frequency,
  next_due_date
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetRecurringPayment :one
SELECT * FROM recurring_payments
WHERE id = $1 LIMIT 1;

-- name: ListRecurringPayments :many
SELECT * FROM recurring_payments
WHERE owner = $1
ORDER BY next_due_date ASC;

-- name: ListDueRecurringPayments :many
SELECT * FROM recurring_payments
WHERE owner = $1
  AND next_due_date <= $2
ORDER BY next_due_date ASC;

-- name: UpdateRecurringPaymentDueDate :one
UPDATE recurring_payments
SET next_due_date = $2
WHERE id = $1
RETURNING *;

-- name: DeleteRecurringPayment :exec
DELETE FROM recurring_payments
WHERE id = $1;
