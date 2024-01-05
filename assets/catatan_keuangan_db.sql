CREATE DATABASE IF NOT EXISTS catatan_keuangan_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');

CREATE TABLE expenses (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    date DATE NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type transaction_type,
    balance DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

INSERT INTO expenses (date, amount, transaction_type, balance, description) VALUES ('2023-12-08', 100, 'CREDIT', 500, 'Pembelian barang');

SELECT * FROM expenses;

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  username VARCHAR(100) unique,
  password VARCHAR(100),
  role VARCHAR(30),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

ALTER TABLE expenses ADD COLUMN user_id uuid;
ALTER TABLE expenses ADD FOREIGN KEY (user_id) REFERENCES users(id);

INSERT INTO users (username, password) VALUES
('john',  'password'),
('doe', 'password'),
('tailor', 'password');