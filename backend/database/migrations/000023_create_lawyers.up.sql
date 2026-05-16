SET NAMES utf8mb4;
CREATE TABLE `lawyers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `photo_url` varchar(512) NOT NULL DEFAULT '',
  `name` varchar(64) NOT NULL DEFAULT '',
  `title` varchar(128) NOT NULL DEFAULT '',
  `tags` varchar(512) NOT NULL DEFAULT '',
  `sort_order` int NOT NULL DEFAULT 0,
  `deleted_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
