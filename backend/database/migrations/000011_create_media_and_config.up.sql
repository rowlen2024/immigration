SET NAMES utf8mb4;
-- Media files table
CREATE TABLE IF NOT EXISTS `media` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `filename` VARCHAR(255) NOT NULL,
    `original_name` VARCHAR(255) NOT NULL DEFAULT '',
    `url` VARCHAR(512) NOT NULL DEFAULT '',
    `mime_type` VARCHAR(64) NOT NULL DEFAULT '',
    `size_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `deleted_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Home page configuration
CREATE TABLE IF NOT EXISTS `home_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(64) NOT NULL,
    `config_value` JSON NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Operation audit log
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `operator_id` BIGINT UNSIGNED NULL,
    `action` VARCHAR(64) NOT NULL,
    `target` VARCHAR(128) NOT NULL DEFAULT '',
    `target_id` BIGINT UNSIGNED NULL,
    `ip` VARCHAR(45) NOT NULL DEFAULT '',
    `details` JSON NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_oplogs_operator` (`operator_id`),
    INDEX `idx_oplogs_action` (`action`),
    INDEX `idx_oplogs_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
