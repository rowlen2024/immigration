CREATE TABLE IF NOT EXISTS `project_categories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(128) NOT NULL,
    `slug` VARCHAR(128) NOT NULL,
    `intro` TEXT NOT NULL,
    `seo_title` VARCHAR(128) NOT NULL DEFAULT '',
    `seo_description` VARCHAR(512) NOT NULL DEFAULT '',
    `status` ENUM('draft', 'published') NOT NULL DEFAULT 'draft',
    `published_at` DATETIME(6) NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_project_categories_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_category_memberships` (
    `project_id` BIGINT UNSIGNED NOT NULL,
    `category_id` BIGINT UNSIGNED NOT NULL,
    `sort_order` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`project_id`, `category_id`),
    KEY `idx_project_category_memberships_category_order` (`category_id`, `sort_order`, `project_id`),
    CONSTRAINT `fk_project_category_memberships_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_category_memberships_category` FOREIGN KEY (`category_id`) REFERENCES `project_categories` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MySQL 会隐式提交 DDL，重试时检测一次性添加的列、索引及外键。
SET @add_navigation_category = IF(
    EXISTS (
        SELECT 1 FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'navigations' AND COLUMN_NAME = 'category_id'
    ),
    'SELECT 1',
    'ALTER TABLE `navigations` ADD COLUMN `category_id` BIGINT UNSIGNED NULL, ADD INDEX `idx_navigations_category_id` (`category_id`), ADD CONSTRAINT `fk_navigations_category` FOREIGN KEY (`category_id`) REFERENCES `project_categories` (`id`) ON DELETE RESTRICT'
);
PREPARE add_navigation_category FROM @add_navigation_category;
EXECUTE add_navigation_category;
DEALLOCATE PREPARE add_navigation_category;
