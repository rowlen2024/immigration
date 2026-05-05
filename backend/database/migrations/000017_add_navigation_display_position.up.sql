SET NAMES utf8mb4;

ALTER TABLE `navigations`
    ADD COLUMN `display_position` VARCHAR(16) NOT NULL DEFAULT 'header' AFTER `status`;

ALTER TABLE `navigations`
    ADD INDEX `idx_nav_display_position` (`display_position`);
