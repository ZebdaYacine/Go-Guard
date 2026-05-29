create database aicha;
use aicha;

create table users(
    ID int8 primary key auto_increment,
    USERNAME varchar(50),
    EMAIL varchar(30),
    PHONE varchar(15),
    PASSWORD varchar(100),
    ROLE INT,
    SEX CHAR(1),
    PICTURE varchar(200)
)