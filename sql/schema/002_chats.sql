-- +goose up
CREATE TABLE chats (
    id VARCHAR(200) NOT NULL PRIMARY KEY,
    createdby VARCHAR(200),
    lastMessage TEXT,
    participants TEXT,
    created_at TIMESTAMP NOT NULL,
    is_group_chat BOOLEAN,
    FOREIGN KEY (createdby) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE(participants, createdby) 
);

-- +goose down
-- DROP TABLE chats;

