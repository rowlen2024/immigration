ALTER TABLE `cases` DROP INDEX `idx_cases_slug`;
ALTER TABLE `cases` DROP COLUMN `slug`;
ALTER TABLE `cases` CHANGE `content` `description` TEXT;
