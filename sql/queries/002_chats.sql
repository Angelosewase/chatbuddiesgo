-- name: CreateChat :execresult
INSERT INTO chats (id, createdby, lastMessage, participants, created_at, is_group_chat)
VALUES (?, ?, ?, ?, ?,?);

-- name: GetChat :one
SELECT * FROM chats
WHERE participants = ?
LIMIT 1;

-- name: UpdateChat :execresult
UPDATE chats
SET lastMessage = ?, participants = ?, is_group_chat = ?
WHERE id = ?;

-- name: UpdateLatestMessage :execresult
UPDATE chats
SET lastMessage = ?
WHERE id = ?;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = ?;

-- name: GetChats :many
SELECT * FROM chats;

-- name: GetChatsByuserId :many
SELECT * FROM chats
WHERE createdby = ?;

-- name: GetUserByName :many
SELECT * FROM users WHERE firstName LIKE CONCAT('%', ?, '%') OR lastName LIKE CONCAT('%', ?, '%');
