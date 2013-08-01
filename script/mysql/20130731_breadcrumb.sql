CREATE TABLE `breadcrumbs` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT, `user_id` BIGINT(20), `lat` DOUBLE(11,6), `lng` DOUBLE(11,6), `acc` DOUBLE(11,6), `timestamp` BIGINT(20), PRIMARY KEY (`id`)) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `breadcrumbs_timestamp` on `breadcrumbs`(`timestamp`);
CREATE INDEX `breadcrumbs_userid_timestamp` on `breadcrumbs`(`user_id`,`timestamp`);

CREATE TABLE `breadcrumbs_windows` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT, `user_id` BIGINT(20), `cross_id` BIGINT(20), `start_at` BIGINT(20), `end_at` BIGINT(20), PRIMARY KEY (`id`)) DEFAULT CHARSET=utf8mb4;
CREATE INDEX `breadcrumbs_windows_userid_crossid_start_end` on `breadcrumbs_windows`(`user_id`,`cross_id`,`start_at`,`end_at`);
CREATE INDEX `breadcrumbs_windows_userid_crossid_end` on `breadcrumbs_windows`(`user_id`,`cross_id`,`end_at`);