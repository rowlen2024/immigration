SET NAMES utf8mb4;

ALTER TABLE `navigations` DROP INDEX `idx_nav_display_position`;
ALTER TABLE `navigations` DROP COLUMN `display_position`;
