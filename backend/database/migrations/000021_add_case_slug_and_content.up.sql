ALTER TABLE `cases` ADD COLUMN `slug` VARCHAR(36) NOT NULL DEFAULT '' AFTER `id`;
UPDATE `cases` SET `slug` = UUID() WHERE `slug` = '';
ALTER TABLE `cases` ADD UNIQUE INDEX `idx_cases_slug` (`slug`);
ALTER TABLE `cases` CHANGE `description` `content` LONGTEXT;
