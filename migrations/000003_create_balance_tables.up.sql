-- migrations/000003_create_balance_tables.up.sql
-- Создание таблиц баланса

CREATE TABLE IF NOT EXISTS transactions (
    user_id VARCHAR(255) NOT NULL,
    t_id VARCHAR(255) PRIMARY KEY,
    amount BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    withdraw BOOLEAN NOT NULL,
    created_at BIGINT NOT NULL
);
-- Базовый индекс для поиска по номеру клиента
CREATE INDEX idx_user_id_t ON transactions(user_id);

CREATE TABLE IF NOT EXISTS balances (
    user_id VARCHAR(255) PRIMARY KEY,
    current_balance BIGINT NOT NULL,
    withdrawn_balance BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
-- Базовый индекс для поиска по номеру клиента
CREATE INDEX idx_user_id_b ON balances(user_id);





