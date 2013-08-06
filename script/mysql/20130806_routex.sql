CREATE TABLE `routex` (`user_id` BIGINT(20) NOT NULL, `cross_id` BIGINT(20) NOT NULL, PRIMARY KEY(`user_id`, `cross_id`), `enable` BOOLEAN, `updated_at` BIGINT(20)) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `routex_cross_id` on `routex`(`cross_id`);
CREATE INDEX `routex_user_id_cross_id` on `routex`(`user_id`, `cross_id`);