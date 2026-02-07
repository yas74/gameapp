-- +migrate Up
Alter TABLE users add column password varchar(255) not null;

-- +migrate Down
Alter TABLE users drop column password;