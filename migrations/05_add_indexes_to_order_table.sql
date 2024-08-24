-- +goose Up
CREATE INDEX order_username_idx ON gofemart.order(username);
CREATE INDEX order_create_date_idx ON gofemart.order(create_date);
CREATE INDEX order_status_idx ON gofemart.order(status);

-- +goose Down
DROP INDEX order_username_idx;
DROP INDEX order_create_date_idx;
DROP INDEX order_status_idx;
