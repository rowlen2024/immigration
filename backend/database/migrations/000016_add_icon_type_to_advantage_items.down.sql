SET NAMES utf8mb4;
-- Remove "icon_type":"lucide" from advantage items.
-- Only touch rows that contain the lucide icon_type.

UPDATE `home_configs`
SET `config_value` = REGEXP_REPLACE(
    `config_value`,
    '"icon_type":"lucide","icon":"([^"]+)"',
    '"icon":"\\1"'
)
WHERE `config_key` = 'advantage_items'
  AND `config_value` LIKE '%"icon_type":"lucide"%';
