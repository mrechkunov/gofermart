-- migrations/000002_create_orders_table.up.sql
-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
    o_number BIGINT PRIMARY KEY,
    o_status VARCHAR(255) NOT NULL,
    o_accrual INT NOT NULL,
    uploaded_at BIGINT NOT NULL,
    created_by VARCHAR(255) NOT NULL
);

-- Базовый индекс для поиска по номеру заказа
CREATE INDEX idx_number ON orders(o_number);


