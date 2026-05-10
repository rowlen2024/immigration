CREATE TABLE IF NOT EXISTS `compare_configs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` BIGINT UNSIGNED NOT NULL,
  `compare_with` JSON NOT NULL,
  `compare_fields` JSON NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_id` (`project_id`),
  CONSTRAINT `fk_cc_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_news` (
  `project_id` BIGINT UNSIGNED NOT NULL,
  `page_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`project_id`, `page_id`),
  CONSTRAINT `fk_pn_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pn_page` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
