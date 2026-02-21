-- +migrate Up
-- Mysql 8.0 set the role value to `user` for all old records
-- dont change the order of the enum values
-- TODO find a better solution than keepping the order!!!
ALTER TABLE `users` ADD COLUMN `role` ENUM('user', 'admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`;