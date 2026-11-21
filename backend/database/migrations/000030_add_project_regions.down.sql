ALTER TABLE `projects` DROP FOREIGN KEY `fk_projects_region`;
ALTER TABLE `projects` DROP INDEX `idx_projects_region_id`, DROP COLUMN `region_id`;
DROP TABLE `project_regions`;
