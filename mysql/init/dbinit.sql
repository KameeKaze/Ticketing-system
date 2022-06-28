SET CHARSET UTF8;
DROP DATABASE IF EXISTS ticketing_system;
CREATE DATABASE ticketing_system DEFAULT CHARACTER SET utf8;

USE ticketing_system;


CREATE TABLE users ( 
    id       INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name     VARCHAR(128) NOT NULL,
    email    VARCHAR(320),
    password VARCHAR(128) NOT NULL,
    role     VARCHAR(32)
);

CREATE TABLE tickets (
    id    INT AUTO_INCREMENT PRIMARY KEY,
    issuer INT,
    FOREIGN KEY (issuer) REFERENCES users(id),
    date DATETIME,
    title TINYTEXT,
    priority INT,
    type VARCHAR(64),
    status INT,
    content TEXT
);


