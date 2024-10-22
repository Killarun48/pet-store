PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT
);

CREATE TABLE IF NOT EXISTS authors
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,    
    name TEXT NOT NULL,
    birth_date DATE FORMAT 'YYYY-MM-DD'
);

CREATE TABLE IF NOT EXISTS books
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    user_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES authors(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO authors (name, birth_date) VALUES
('test','2000-01-01'),
('test2','2005-01-02');

INSERT INTO users (username, email) VALUES
('test','test@t.com'),
('test2','test2@t.com');

INSERT INTO books (title, author_id, user_id) VALUES
('test2',2,2),
('test3',2,2);