SET NAMES utf8mb4;

ALTER TABLE `users` MODIFY COLUMN `role` VARCHAR(64) NOT NULL DEFAULT 'viewer';

CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(64) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `description` VARCHAR(255) NOT NULL DEFAULT '',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=disabled',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(96) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `module` VARCHAR(64) NOT NULL,
    `action` VARCHAR(32) NOT NULL,
    `sort_order` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permissions_code` (`code`),
    KEY `idx_permissions_module` (`module`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `role_permissions` (
    `role_id` BIGINT UNSIGNED NOT NULL,
    `permission_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`role_id`, `permission_id`),
    CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_permission_overrides` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `permission_id` BIGINT UNSIGNED NOT NULL,
    `effect` ENUM('allow','deny') NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_permission` (`user_id`, `permission_id`),
    CONSTRAINT `fk_user_permission_overrides_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_permission_overrides_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO `roles` (`code`, `name`, `description`, `status`, `is_system`) VALUES
('admin', '管理员', '拥有后台全部权限', 1, 1),
('editor', '编辑者', '负责内容维护和线索查看', 1, 1),
('viewer', '只读用户', '仅可查看后台内容', 1, 1);

INSERT IGNORE INTO `permissions` (`code`, `name`, `module`, `action`, `sort_order`) VALUES
('dashboard:read', '查看控制台', 'dashboard', 'read', 10),
('projects:read', '查看项目', 'projects', 'read', 20),
('projects:write', '维护项目', 'projects', 'write', 21),
('homepage:read', '查看首页配置', 'homepage', 'read', 30),
('homepage:write', '维护首页配置', 'homepage', 'write', 31),
('navigation:read', '查看导航', 'navigation', 'read', 40),
('navigation:write', '维护导航', 'navigation', 'write', 41),
('pages:read', '查看页面', 'pages', 'read', 50),
('pages:write', '维护页面', 'pages', 'write', 51),
('media:read', '查看媒体库', 'media', 'read', 60),
('media:write', '维护媒体库', 'media', 'write', 61),
('faqs:read', '查看 FAQ', 'faqs', 'read', 70),
('faqs:write', '维护 FAQ', 'faqs', 'write', 71),
('cases:read', '查看案例', 'cases', 'read', 80),
('cases:write', '维护案例', 'cases', 'write', 81),
('lawyers:read', '查看律师团队', 'lawyers', 'read', 90),
('lawyers:write', '维护律师团队', 'lawyers', 'write', 91),
('testimonials:read', '查看客户评价', 'testimonials', 'read', 100),
('testimonials:write', '维护客户评价', 'testimonials', 'write', 101),
('leads:read', '查看咨询', 'leads', 'read', 110),
('leads:write', '维护咨询', 'leads', 'write', 111),
('settings:read', '查看网站设置', 'settings', 'read', 120),
('settings:write', '维护网站设置', 'settings', 'write', 121),
('users:read', '查看用户', 'users', 'read', 130),
('users:write', '维护用户', 'users', 'write', 131),
('roles:read', '查看角色权限', 'roles', 'read', 140),
('roles:write', '维护角色权限', 'roles', 'write', 141);

INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM `roles` r JOIN `permissions` p
WHERE r.code = 'admin';

INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM `roles` r JOIN `permissions` p
WHERE r.code = 'editor'
  AND p.code IN (
    'dashboard:read',
    'projects:read',
    'homepage:read','homepage:write',
    'navigation:read','navigation:write',
    'pages:read','pages:write',
    'media:read','media:write',
    'faqs:read','faqs:write',
    'cases:read','cases:write',
    'lawyers:read','lawyers:write',
    'testimonials:read','testimonials:write',
    'leads:read',
    'settings:read'
  );

INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM `roles` r JOIN `permissions` p
WHERE r.code = 'viewer'
  AND p.code IN (
    'dashboard:read',
    'projects:read',
    'homepage:read',
    'navigation:read',
    'pages:read',
    'media:read',
    'faqs:read',
    'cases:read',
    'lawyers:read',
    'testimonials:read',
    'leads:read',
    'settings:read'
  );
