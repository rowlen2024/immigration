SET NAMES utf8mb4;
CREATE TABLE IF NOT EXISTS `navigations` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `label` VARCHAR(255) NOT NULL,
    `link` VARCHAR(512) NOT NULL,
    `parent_id` BIGINT UNSIGNED NULL,
    `sort_order` INT NOT NULL DEFAULT 0,
    `status` TINYINT(1) NOT NULL DEFAULT 1,
    `deleted_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_navigations_parent` (`parent_id`),
    INDEX `idx_navigations_sort` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `navigations` (`id`, `label`, `link`, `parent_id`, `sort_order`, `status`) VALUES
(1, '美国EB-5', '/usa/eb5', NULL, 1, 1),
(2, '香港投资', '/hongkong/cies', NULL, 2, 1),
(3, '巴拿马购房', '/panama/property', NULL, 3, 1),
(4, '项目对比', '/compare', NULL, 4, 1),
(5, '关于我们', '/about', NULL, 5, 1),
(18, '公司简介', '/about', 5, 1, 1),
(19, '成功案例', '/cases', 5, 2, 1),
(20, '常见问题', '/faq', 5, 3, 1),
(21, '联系我们', '/contact', 5, 4, 1);
