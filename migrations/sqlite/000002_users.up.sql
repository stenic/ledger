
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(50) NOT NULL,
	password VARCHAR(200) NOT NULL,
	last_login TIMESTAMP
);
