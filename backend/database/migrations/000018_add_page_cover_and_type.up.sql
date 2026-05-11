ALTER TABLE pages ADD COLUMN `cover_image` VARCHAR(512) NOT NULL DEFAULT '' AFTER `content`;
ALTER TABLE pages ADD COLUMN `page_type` VARCHAR(32) NOT NULL DEFAULT 'default' AFTER `template`;
