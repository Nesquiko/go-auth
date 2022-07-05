INSERT INTO users (username, email, passwordHash, secret2FA, enabled2FA)
VALUES ("nesquiko","nesquiko@foo.com", "$2a$10$NCrqADHPMllaWXxmpqvUA.6q0NFenzjo4vjjb/289F5wrQnyvhPGm", NULL, 0);


UPDATE users SET secret2FA = 'EESQTH3G2YF26LUF' WHERE username = 'nesquiko';

SELECT * FROM users;