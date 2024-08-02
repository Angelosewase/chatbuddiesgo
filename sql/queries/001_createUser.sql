-- name: CreateUser :execresult
INSERT INTO users (id, firstName, lastName, email, password)
VALUES (?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;