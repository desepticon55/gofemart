-- +goose Up
CREATE TABLE gofemart.balance
(
    username VARCHAR(255) UNIQUE NOT NULL,
    balance  NUMERIC(18, 2)      NOT NULL,
    opt_lock bigint              NOT NULL,
    PRIMARY KEY (username)
);

CREATE TABLE gofemart.withdrawal
(
    id           UUID UNIQUE    NOT NULL,
    order_number VARCHAR(255)   NOT NULL,
    username     VARCHAR(255)   NOT NULL,
    sum          NUMERIC(18, 2) NOT NULL,
    create_date  TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE gofemart.balance;
DROP TABLE gofemart.withdrawal;