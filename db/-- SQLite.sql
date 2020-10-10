-- SQLite

PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：users
CREATE TABLE users (id  INTEGER,name  VARCHAR NOT NULL,password  VARCHAR NOT NULL,token  VARCHAR,expire_time  INTEGER DEFAULT 0,fail_num  INTEGER DEFAULT 0,fail_time  INTEGER DEFAULT 0,PRIMARY KEY (id ASC),UNIQUE (name ASC));
INSERT INTO users (id, name, password, token, expire_time, fail_num, fail_time) VALUES (1, 'admin', '$2a$10$xgCuBT7Fphk1EHFcLh8PXeZ4w.gg2XAivq9tTm6iMqc8lWs8sP3Yy', '', 0, 0, 0);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
