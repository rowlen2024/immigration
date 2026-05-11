CREATE TABLE IF NOT EXISTS `project_advantages` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `icon` VARCHAR(64) NOT NULL DEFAULT '',
    `icon_type` VARCHAR(32) NOT NULL DEFAULT 'lucide',
    `title` VARCHAR(128) NOT NULL,
    `description` TEXT,
    `sort_order` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME,
    PRIMARY KEY (`id`),
    KEY `idx_project_advantages_project_id` (`project_id`),
    CONSTRAINT `fk_project_advantages_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
