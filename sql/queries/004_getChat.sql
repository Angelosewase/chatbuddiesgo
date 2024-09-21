-- name: GetChatByChatId :one
SELECT * FROM chats WHERE id= ? ;