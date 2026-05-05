SET NAMES utf8mb4;
-- Add "icon_type":"lucide" to each advantage item that lacks it.
-- REGEXP_REPLACE inserts icon_type before every "icon" field.
-- The WHERE guard ensures this is idempotent.

UPDATE `home_configs`
SET `config_value` = REGEXP_REPLACE(
    `config_value`,
    '"icon":"([^"]+)"',
    '"icon_type":"lucide","icon":"\\1"'
)
WHERE `config_key` = 'advantage_items'
  AND `config_value` NOT LIKE '%"icon_type"%';
