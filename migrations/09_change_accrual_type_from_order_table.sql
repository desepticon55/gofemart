-- +goose Up
ALTER TABLE gofemart.order ALTER COLUMN accrual TYPE NUMERIC(18, 2);

-- +goose Down
ALTER TABLE gofemart.order ALTER COLUMN accrual TYPE BIGINT;
