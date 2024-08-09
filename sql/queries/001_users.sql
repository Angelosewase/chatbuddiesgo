-- name: CreateUser :execresult
INSERT INTO users (id, firstName, lastName, email, password)
VALUES (?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: UpdateUser :execresult
UPDATE users
SET firstName = ?, lastName = ?, email = ?, password = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetUsers :many
SELECT * FROM users;
