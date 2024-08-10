-- +goose Up
CREATE TABLE gofemart.user
(
    username VARCHAR(255) UNIQUE NOT NULL,
    password TEXT                NOT NULL,
    PRIMARY KEY (username)
);

-- +goose Down
DROP TABLE gofemart.user;