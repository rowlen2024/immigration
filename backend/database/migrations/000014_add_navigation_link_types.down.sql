-- Remove foreign key constraints first
ALTER TABLE `navigations`
    DROP FOREIGN KEY `fk_navigations_project`,
    DROP FOREIGN KEY `fk_navigations_page`;

-- Remove indices
ALTER TABLE `navigations`
    DROP INDEX `idx_nav_link_type`,
    DROP INDEX `idx_nav_project_id`,
    DROP INDEX `idx_nav_page_id`;

-- Restore link to NOT NULL (ensure no null links exist first)
UPDATE `navigations` SET `link` = '' WHERE `link` IS NULL;
ALTER TABLE `navigations`
    MODIFY COLUMN `link` VARCHAR(512) NOT NULL;

-- Drop the added columns
ALTER TABLE `navigations`
    DROP COLUMN `page_id`,
    DROP COLUMN `project_id`,
    DROP COLUMN `link_type`;
