-- migrations/000001_create_users_table.up.sql
-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    ulogin VARCHAR(255) PRIMARY KEY,
    upassword VARCHAR(255) NOT NULL,
    uBearer VARCHAR(255) NOT NULL
);

-- Базовый индекс для поиска по логину
CREATE INDEX idx_login ON users(ulogin);