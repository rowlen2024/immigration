ALTER TABLE `navigations`
    DROP FOREIGN KEY `fk_navigations_category`,
    DROP INDEX `idx_navigations_category_id`,
    DROP COLUMN `category_id`;

DROP TABLE `project_category_memberships`;
DROP TABLE `project_categories`;
