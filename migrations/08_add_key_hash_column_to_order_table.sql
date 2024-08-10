-- +goose Up
ALTER TABLE gofemart.order ADD COLUMN key_hash BIGINT;
ALTER TABLE gofemart.order ADD COLUMN key_hash_module INT;

-- +goose Down
ALTER TABLE gofemart.order DROP COLUMN key_hash;
ALTER TABLE gofemart.order DROP COLUMN key_hash_module;
