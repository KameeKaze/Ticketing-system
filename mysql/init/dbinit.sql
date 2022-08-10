SET CHARSET UTF8;
DROP DATABASE IF EXISTS ticketing_system;
CREATE DATABASE ticketing_system DEFAULT CHARACTER SET utf8;

USE ticketing_system;


CREATE TABLE users ( 
    id       VARCHAR(36)  NOT NULL PRIMARY KEY,
    name     VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    role     VARCHAR(32)  NOT NULL
);

CREATE TABLE tickets (
    id       VARCHAR(36) NOT NULL PRIMARY KEY,
    issuer   VARCHAR(36) NOT NULL,
    date     DATETIME    NOT NULL,
    title    TINYTEXT    NOT NULL,
    status   INT         NOT NULL,
    content  TEXT        NOT NULL,
    FOREIGN KEY (issuer) REFERENCES users(id)
);

-- admin:admin
INSERT INTO users (id, name, password, role) VALUES (UUID(), "admin", "$2a$07$qG7WMfWkMvsTgQfCEaMxtex5VOYIZHlZX0cTRnsW0tH3gE.br4xKK", "admin");
