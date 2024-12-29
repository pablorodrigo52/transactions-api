CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    transaction_date TEXT NOT NULL,
    purchase_amount REAL NOT NULL
);