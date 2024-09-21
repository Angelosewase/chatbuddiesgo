-- +goose Up
CREATE TABLE messages (
    id VARCHAR(200) PRIMARY KEY,
    chat_id VARCHAR(200) NOT NULL,
    sender_id  VARCHAR(200)  NOT NULL,
    content TEXT NOT NULL,
    content_type ENUM('text', 'image', 'video', 'file') DEFAULT 'text',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose down 
-- DROP TABLE messages;