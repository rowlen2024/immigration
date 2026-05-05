SET NAMES utf8mb4;

-- Add link_type and reference columns to navigations
ALTER TABLE `navigations`
    ADD COLUMN `link_type` VARCHAR(32) NOT NULL DEFAULT 'custom' AFTER `link`,
    ADD COLUMN `project_id` BIGINT UNSIGNED NULL AFTER `link_type`,
    ADD COLUMN `page_id` BIGINT UNSIGNED NULL AFTER `project_id`;

-- Make link column nullable for project/page link types (link is generated dynamically)
ALTER TABLE `navigations`
    MODIFY COLUMN `link` VARCHAR(512) NULL;

-- Add foreign key constraints
ALTER TABLE `navigations`
    ADD CONSTRAINT `fk_navigations_project` FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE RESTRICT,
    ADD CONSTRAINT `fk_navigations_page` FOREIGN KEY (`page_id`) REFERENCES `pages`(`id`) ON DELETE RESTRICT;

-- Add indexes for reference lookups
ALTER TABLE `navigations`
    ADD INDEX `idx_nav_link_type` (`link_type`),
    ADD INDEX `idx_nav_project_id` (`project_id`),
    ADD INDEX `idx_nav_page_id` (`page_id`);

-- Attempt to auto-match existing navigations to projects by slug.
-- Only match entries whose link is exactly /projects/<slug> or /
-- (exclude hash-fragment child links like /usa/eb5#requirements).
UPDATE `navigations` n
    JOIN `projects` p ON n.link LIKE CONCAT('%', p.slug, '%') AND n.link NOT LIKE '%#%'
SET n.link_type = 'project', n.project_id = p.id, n.link = CONCAT('/projects/', p.slug)
WHERE n.link_type = 'custom' AND n.link NOT LIKE '%#%';

-- Attempt to auto-match existing navigations to pages by slug.
-- Pages use route /pages/:slug (or /:slug for simple paths like /about, /contact).
UPDATE `navigations` n
    JOIN `pages` pg ON n.link = CONCAT('/', pg.slug)
SET n.link_type = 'page', n.page_id = pg.id, n.link = CONCAT('/pages/', pg.slug)
WHERE n.link_type = 'custom';

-- Try matching pages with nested slugs (pages slug can contain / for nested paths)
UPDATE `navigations` n
    JOIN `pages` pg ON n.link = CONCAT('/pages/', pg.slug)
SET n.link_type = 'page', n.page_id = pg.id
WHERE n.link_type = 'custom';
