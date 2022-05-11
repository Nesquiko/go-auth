INSERT INTO users (username, email, passwordHash)
VALUES ('nesquiko', 'nesquiko@foo.com', '$2a$10$NCrqADHPMllaWXxmpqvUA.6q0NFenzjo4vjjb/289F5wrQnyvhPGm'),
        ('vava', 'vava@bar.sk', '$2a$10$3rQ.UC9d95DttaJ5yRBXJOh/SMUATcfprGcDSfbxGgWlkYvOm3NvC');

SELECT * FROM users;