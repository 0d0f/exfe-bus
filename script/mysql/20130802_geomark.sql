CREATE TABLE `geomarks` (`id` VARCHAR(20) NOT NULL, `type` VARCHAR(20) NOT NULL, `cross_id` BIGINT(20) NOT NULL, `mark` TEXT, `touched_at` BIGINT, `deleted` BOOLEAN) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `geomarks_id_cross_id` on `geomarks`(`id`, `type`, `cross_id`, `deleted`);
CREATE INDEX `geomarks_cross_id` on `geomarks`(`cross_id`, `deleted`);