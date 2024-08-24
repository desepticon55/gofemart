-- +goose Up
ALTER TABLE gofemart.order ADD COLUMN accrual BIGINT;

-- +goose Down
ALTER TABLE gofemart.order DROP COLUMN accrual;
