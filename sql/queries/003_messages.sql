-- name: AddTextMessage :execresult
INSERT INTO messages(id,chat_id,sender_id,content) VALUES (?,?,?,?);

-- name: GetMessagesByChatId :many
SELECT * FROM messages WHERE chat_id = ? ;

-- name: GetMessageById :one 
SELECT * FROM messages WHERE id = ?;

-- name: DeleteMessage :execresult
UPDATE messages SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ?;
