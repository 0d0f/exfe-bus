CREATE TABLE `geomarks` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT, PRIMARY KEY(`id`), `cross_id` BIGINT(20) NOT NULL, `mark` TEXT, `touched_at` BIGINT) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `geomarks_id_cross_id` on `geomarks`(`id`, `cross_id`);
CREATE INDEX `geomarks_cross_id` on `geomarks`(`cross_id`);