SET NAMES utf8mb4;
CREATE TABLE IF NOT EXISTS `cost_items` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `amount` VARCHAR(64) NOT NULL DEFAULT '',
    `amount_value` DECIMAL(12,2) NULL,
    `amount_currency` CHAR(3) NOT NULL DEFAULT 'USD',
    `note` VARCHAR(512) NOT NULL DEFAULT '',
    `sort_order` INT NOT NULL DEFAULT 0,
    `deleted_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_cost_items_project` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
