SET NAMES utf8mb4;
CREATE TABLE IF NOT EXISTS `cases` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NULL,
    `name` VARCHAR(128) NOT NULL DEFAULT '',
    `country_from` VARCHAR(64) NOT NULL DEFAULT '',
    `investment_amount` VARCHAR(64) NOT NULL DEFAULT '',
    `investment_value` DECIMAL(12,2) NULL,
    `processing_period` VARCHAR(64) NOT NULL DEFAULT '',
    `description` TEXT,
    `photo_url` VARCHAR(512) NOT NULL DEFAULT '',
    `sort_order` INT NOT NULL DEFAULT 0,
    `deleted_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_cases_project` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
