CREATE TABLE IF NOT EXISTS versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    environment VARCHAR(100),
    application VARCHAR(200),
    version VARCHAR(100),
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
