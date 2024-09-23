-- name: GetChatNotCreatedByTheUser :many

SELECT * FROM chats WHERE createdby != ? AND participants LIKE CONCAT('%',? , '%');