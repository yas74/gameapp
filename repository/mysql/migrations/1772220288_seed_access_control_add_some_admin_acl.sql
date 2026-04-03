-- +migrate Up

INSERT INTO `permissions` (`id`, `title`) VALUES(1, 'user-list');
INSERT INTO `permissions` (`id`, `title`) VALUES(2, 'user-delete');

INSERT INTO `access_controls` (`id`, `actor_id`, `actor_type`, `permission_id`) VALUES(1, 2, 'role', 1);
INSERT INTO `access_controls` (`id`, `actor_id`, `actor_type`, `permission_id`) VALUES(2, 2, 'role', 2);

-- +migrate Down

DELETE FROM `permissions` WHERE id IN (1,2);

DELETE FROM `access_controls` WHERE id IN (1,2);