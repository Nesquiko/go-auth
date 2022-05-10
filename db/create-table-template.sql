DROP TABLE IF EXISTS users;
CREATE TABLE users(
    uuid VARCHAR(36) DEFAULT (uuid()) NOT NULL PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL UNIQUE,
    passwordHash CHAR(60) BINARY NOT NULL
);

INSERT INTO users (username, email, passwordHash)
VALUES ('nesquiko', 'nesquiko@foo.com', '$2a$10$NCrqADHPMllaWXxmpqvUA.6q0NFenzjo4vjjb/289F5wrQnyvhPGm');