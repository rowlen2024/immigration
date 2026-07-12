ALTER TABLE `cases`
    ADD COLUMN `is_pinned` TINYINT(1) NOT NULL DEFAULT 0 AFTER `sort_order`;
