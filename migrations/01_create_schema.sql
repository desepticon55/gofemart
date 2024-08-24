-- +goose Up
CREATE SCHEMA gofemart AUTHORIZATION postgres;
GRANT USAGE ON SCHEMA gofemart TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gofemart TO postgres;

-- +goose Down
DROP SCHEMA gofemart;