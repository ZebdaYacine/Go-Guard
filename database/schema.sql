create database aicha;
use aicha;

create table users(
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username varchar(50) UNIQUE,
    email varchar(30) UNIQUE,
    phone varchar(40) UNIQUE,
    password varchar(100),
    role INT,
    sex CHAR(1),
    picture varchar(200)
)