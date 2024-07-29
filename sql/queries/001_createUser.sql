-- name: CreateUser :execresult

INSERT INTO users (id,firstName,lastName,email,password)VALUES(
    ?,?,?,?,?
)