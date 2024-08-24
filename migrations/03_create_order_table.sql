-- +goose Up
CREATE TABLE gofemart.order
(
    order_number     VARCHAR(255) UNIQUE      NOT NULL,
    username         VARCHAR(255)             NOT NULL,
    create_date      TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modify_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status           VARCHAR(50)              NOT NULL,
    PRIMARY KEY (order_number)
);

-- +goose Down
DROP TABLE gofemart.order;