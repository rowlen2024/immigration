SET NAMES utf8mb4;
SET SESSION group_concat_max_len = 1048576;

-- 临时保存按原配置顺序转换后的项目 ID。
CREATE TEMPORARY TABLE `tmp_project_showcase_ids` (
    `config_id` BIGINT UNSIGNED NOT NULL,
    `featured_project_ids` JSON NOT NULL,
    PRIMARY KEY (`config_id`)
);

INSERT INTO `tmp_project_showcase_ids` (`config_id`, `featured_project_ids`)
SELECT
    matched.`config_id`,
    COALESCE(
        JSON_EXTRACT(
            CONCAT(
                '[',
                GROUP_CONCAT(matched.`project_id` ORDER BY matched.`first_position` SEPARATOR ','),
                ']'
            ),
            '$'
        ),
        JSON_ARRAY()
    )
FROM (
    SELECT
        config.`id` AS `config_id`,
        project.`id` AS `project_id`,
        MIN(item.`position`) AS `first_position`
    FROM `home_configs` AS config
    LEFT JOIN JSON_TABLE(
        COALESCE(JSON_EXTRACT(config.`config_value`, '$.featured_slugs'), JSON_ARRAY()),
        '$[*]' COLUMNS (
            `position` FOR ORDINALITY,
            `slug` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci PATH '$'
        )
    ) AS item ON TRUE
    LEFT JOIN `projects` AS project
        ON project.`slug` = item.`slug`
        AND project.`deleted_at` IS NULL
    WHERE config.`config_key` = 'project_showcase'
    GROUP BY config.`id`, project.`id`
) AS matched
GROUP BY matched.`config_id`;

UPDATE `home_configs` AS config
INNER JOIN `tmp_project_showcase_ids` AS migrated
    ON migrated.`config_id` = config.`id`
SET config.`config_value` = JSON_SET(
    JSON_REMOVE(
        config.`config_value`,
        '$.featured_slugs',
        '$.featured_projects'
    ),
    '$.featured_project_ids',
    migrated.`featured_project_ids`
)
WHERE config.`config_key` = 'project_showcase';

DROP TEMPORARY TABLE `tmp_project_showcase_ids`;
