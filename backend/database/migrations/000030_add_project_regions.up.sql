SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `project_regions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` CHAR(2) NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `flag_emoji` VARCHAR(16) NOT NULL DEFAULT '',
    `sort_order` INT NOT NULL DEFAULT 0,
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_project_regions_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MySQL 会隐式提交 DDL，启动迁移重试时跳过已执行的原子 ALTER。
SET @add_project_region = IF(
    EXISTS (
        SELECT 1 FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'projects' AND COLUMN_NAME = 'region_id'
    ),
    'SELECT 1',
    'ALTER TABLE `projects` ADD COLUMN `region_id` BIGINT UNSIGNED NULL AFTER `country`, ADD INDEX `idx_projects_region_id` (`region_id`), ADD CONSTRAINT `fk_projects_region` FOREIGN KEY (`region_id`) REFERENCES `project_regions` (`id`) ON DELETE RESTRICT'
);
PREPARE add_project_region FROM @add_project_region;
EXECUTE add_project_region;
DEALLOCATE PREPARE add_project_region;

INSERT INTO `project_regions` (`code`, `name`, `flag_emoji`, `sort_order`) VALUES
    ('US', '美国', '🇺🇸', 10),
    ('PA', '巴拿马', '🇵🇦', 20),
    ('HK', '中国香港', '🇭🇰', 30),
    ('AE', '阿联酋', '🇦🇪', 100),
    ('AG', '安提瓜和巴布达', '🇦🇬', 100),
    ('CN', '中国', '🇨🇳', 100),
    ('DM', '多米尼克', '🇩🇲', 100),
    ('ES', '西班牙', '🇪🇸', 100),
    ('GB', '英国', '🇬🇧', 100),
    ('GD', '格林纳达', '🇬🇩', 100),
    ('GR', '希腊', '🇬🇷', 100),
    ('HU', '匈牙利', '🇭🇺', 100),
    ('JP', '日本', '🇯🇵', 100),
    ('KN', '圣基茨和尼维斯', '🇰🇳', 100),
    ('LC', '圣卢西亚', '🇱🇨', 100),
    ('MT', '马耳他', '🇲🇹', 100),
    ('MY', '马来西亚', '🇲🇾', 100),
    ('PT', '葡萄牙', '🇵🇹', 100),
    ('SG', '新加坡', '🇸🇬', 100),
    ('ST', '圣多美和普林西比', '🇸🇹', 100),
    ('TH', '泰国', '🇹🇭', 100),
    ('TR', '土耳其', '🇹🇷', 100),
    ('VU', '瓦努阿图', '🇻🇺', 100)
ON DUPLICATE KEY UPDATE `code` = `project_regions`.`code`;

-- 旧字段保留原样。地区文案与项目明显不一致的存量记录留给运营确认。
UPDATE `projects` AS p
JOIN `project_regions` AS r
    ON CAST(TRIM(p.`country`) AS BINARY) = CAST(r.`name` AS BINARY)
    OR (r.`code` = 'HK' AND CAST(TRIM(p.`country`) AS BINARY) = CAST('中国（香港特别行政区）' AS BINARY))
SET p.`region_id` = r.`id`, p.`updated_at` = p.`updated_at`
WHERE p.`deleted_at` IS NULL
    AND p.`region_id` IS NULL
    AND NOT (p.`slug` = 'xiang-gang-gao-duan-ren-cai' AND r.`code` = 'CN');
