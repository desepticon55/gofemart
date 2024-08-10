-- +goose Up
ALTER TABLE gofemart.order ADD COLUMN opt_lock BIGINT;

-- +goose Down
ALTER TABLE gofemart.order DROP COLUMN opt_lock;
