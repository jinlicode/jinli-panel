-- SQLite

PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：users
CREATE TABLE users (id INTEGER PRIMARY KEY, name VARCHAR UNIQUE NOT NULL, password VARCHAR NOT NULL);
INSERT INTO users (id, name, password) VALUES (1, 'admin', '$2a$10$zGn7.ytKFWpjYI6mRZjaQ.lIf0PwIVRZZLfGHEtHCVTzdvWHuQQ5q');

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
