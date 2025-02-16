-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    coins INTEGER DEFAULT 1000,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Создание таблицы сессий
CREATE TABLE IF NOT EXISTS sessions (
    user_id INTEGER REFERENCES users(id),
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- Создание таблицы инвентаря
CREATE TABLE IF NOT EXISTS inventory (
    user_id INTEGER REFERENCES users(id),
    type VARCHAR(255) NOT NULL,
    quantity INTEGER DEFAULT 0,
    PRIMARY KEY (user_id, type)
);

-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS coin_transactions (
    id SERIAL PRIMARY KEY,
    from_user INTEGER REFERENCES users(id),
    to_user INTEGER REFERENCES users(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);