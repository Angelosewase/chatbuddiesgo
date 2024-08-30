-- +goose up
CREATE TABLE users(
    id VARCHAR(200) PRIMARY KEY ,
    firstName TEXT,
    lastName TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

-- +goose down
DROP TABLE users;