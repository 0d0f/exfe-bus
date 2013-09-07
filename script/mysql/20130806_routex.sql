CREATE TABLE `routex` (`cross_id` BIGINT(20) NOT NULL, PRIMARY KEY(`cross_id`), `updated_at` BIGINT(20)) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `routex_cross_id` on `routex`(`cross_id`);
CREATE INDEX `routex_cross_id_updated_at` on `routex`(`cross_id`, `updated_at`);