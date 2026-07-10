package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const projectShowcaseMigrationName = "000025_migrate_project_showcase_to_project_ids.up.sql"

func TestProjectShowcaseMigration(t *testing.T) {
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("未设置 MIGRATION_TEST_DSN，跳过迁移集成测试")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接迁移测试数据库失败: %v", err)
	}
	var databaseName string
	if err := db.Raw("SELECT DATABASE()").Scan(&databaseName).Error; err != nil {
		t.Fatalf("获取迁移测试数据库名失败: %v", err)
	}
	if !strings.HasSuffix(databaseName, "_migration_test") {
		t.Fatalf("迁移测试只能连接名称以 _migration_test 结尾的数据库，当前为 %q", databaseName)
	}

	migrationsDir := writeProjectShowcaseMigration(t)
	t.Cleanup(func() { resetProjectShowcaseMigrationTables(t, db) })

	t.Run("保留顺序并清理重复和失效引用", func(t *testing.T) {
		prepareProjectShowcaseMigrationTables(t, db)
		for _, statement := range []string{
			"INSERT INTO projects (id, slug, deleted_at) VALUES (1, 'one', NULL), (2, 'two', NULL), (3, 'three', NULL), (4, 'deleted', NOW())",
			`INSERT INTO home_configs (id, config_key, config_value) VALUES (1, 'project_showcase', '{"section_title":"精选项目","section_subtitle":"副标题","featured_slugs":["three","one","two","one","missing","deleted"],"featured_projects":[{"slug":"stale"}],"extra":{"keep":true}}')`,
		} {
			if err := db.Exec(statement).Error; err != nil {
				t.Fatalf("准备迁移数据失败: %v", err)
			}
		}

		if err := RunMigrations(db, migrationsDir); err != nil {
			t.Fatalf("执行迁移失败: %v", err)
		}
		config := readProjectShowcaseConfig(t, db)
		assertProjectIDs(t, config, []uint64{3, 1, 2})
		assertMissingJSONField(t, config, "featured_slugs")
		assertMissingJSONField(t, config, "featured_projects")
		if string(config["section_title"]) != `"精选项目"` || string(config["section_subtitle"]) != `"副标题"` {
			t.Fatalf("迁移未保留标题字段: %s", mustMarshalJSON(t, config))
		}
		if string(config["extra"]) != `{"keep": true}` && string(config["extra"]) != `{"keep":true}` {
			t.Fatalf("迁移未保留额外字段: %s", mustMarshalJSON(t, config))
		}
	})

	t.Run("空数组保持为空且删除旧字段", func(t *testing.T) {
		prepareProjectShowcaseMigrationTables(t, db)
		if err := db.Exec(`INSERT INTO home_configs (id, config_key, config_value) VALUES (1, 'project_showcase', '{"featured_slugs":[],"featured_projects":[]}')`).Error; err != nil {
			t.Fatalf("准备迁移数据失败: %v", err)
		}

		if err := RunMigrations(db, migrationsDir); err != nil {
			t.Fatalf("执行迁移失败: %v", err)
		}
		config := readProjectShowcaseConfig(t, db)
		assertProjectIDs(t, config, []uint64{})
		assertMissingJSONField(t, config, "featured_slugs")
		assertMissingJSONField(t, config, "featured_projects")
	})

	t.Run("长数组不会被截断", func(t *testing.T) {
		prepareProjectShowcaseMigrationTables(t, db)
		const projectCount = 350
		slugs := make([]string, 0, projectCount)
		expectedIDs := make([]uint64, 0, projectCount)
		for i := projectCount; i >= 1; i-- {
			slug := fmt.Sprintf("project-%03d", i)
			if err := db.Exec("INSERT INTO projects (id, slug, deleted_at) VALUES (?, ?, NULL)", i, slug).Error; err != nil {
				t.Fatalf("插入长数组项目失败: %v", err)
			}
			slugs = append(slugs, slug)
			expectedIDs = append(expectedIDs, uint64(i))
		}
		configValue, err := json.Marshal(map[string]any{"featured_slugs": slugs})
		if err != nil {
			t.Fatalf("构造长数组配置失败: %v", err)
		}
		if err := db.Exec("INSERT INTO home_configs (id, config_key, config_value) VALUES (?, ?, ?)", 1, "project_showcase", configValue).Error; err != nil {
			t.Fatalf("插入长数组配置失败: %v", err)
		}

		if err := RunMigrations(db, migrationsDir); err != nil {
			t.Fatalf("执行迁移失败: %v", err)
		}
		assertProjectIDs(t, readProjectShowcaseConfig(t, db), expectedIDs)
	})

	t.Run("没有配置行时直接跳过", func(t *testing.T) {
		prepareProjectShowcaseMigrationTables(t, db)
		if err := RunMigrations(db, migrationsDir); err != nil {
			t.Fatalf("执行迁移失败: %v", err)
		}
		var count int64
		if err := db.Model(&struct{}{}).Table("home_configs").Count(&count).Error; err != nil {
			t.Fatalf("统计首页配置失败: %v", err)
		}
		if count != 0 {
			t.Fatalf("没有项目展示配置时不应新增数据，实际新增 %d 行", count)
		}
	})
}

func writeProjectShowcaseMigration(t *testing.T) string {
	t.Helper()
	sourcePath := filepath.Join("..", "..", "database", "migrations", projectShowcaseMigrationName)
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("读取迁移文件失败: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, projectShowcaseMigrationName), content, 0o600); err != nil {
		t.Fatalf("准备迁移测试目录失败: %v", err)
	}
	return dir
}

func prepareProjectShowcaseMigrationTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	resetProjectShowcaseMigrationTables(t, db)
	statements := []string{
		"CREATE TABLE projects (id BIGINT UNSIGNED NOT NULL PRIMARY KEY, slug VARCHAR(255) NOT NULL UNIQUE, deleted_at DATETIME NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"CREATE TABLE home_configs (id BIGINT UNSIGNED NOT NULL PRIMARY KEY, config_key VARCHAR(100) NOT NULL UNIQUE, config_value JSON NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
	}
	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("创建迁移测试表失败: %v", err)
		}
	}
}

func resetProjectShowcaseMigrationTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, table := range []string{"migrations", "home_configs", "projects"} {
		if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			t.Fatalf("清理迁移测试表 %s 失败: %v", table, err)
		}
	}
}

func readProjectShowcaseConfig(t *testing.T, db *gorm.DB) map[string]json.RawMessage {
	t.Helper()
	var value string
	if err := db.Raw("SELECT config_value FROM home_configs WHERE config_key = ?", "project_showcase").Scan(&value).Error; err != nil {
		t.Fatalf("读取迁移结果失败: %v", err)
	}
	var config map[string]json.RawMessage
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		t.Fatalf("解析迁移结果失败: %v", err)
	}
	return config
}

func assertProjectIDs(t *testing.T, config map[string]json.RawMessage, expected []uint64) {
	t.Helper()
	var actual []uint64
	if err := json.Unmarshal(config["featured_project_ids"], &actual); err != nil {
		t.Fatalf("解析 featured_project_ids 失败: %v", err)
	}
	if len(actual) != len(expected) {
		t.Fatalf("featured_project_ids 长度不一致，期望 %d，实际 %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("featured_project_ids[%d] 期望 %d，实际 %d", i, expected[i], actual[i])
		}
	}
}

func assertMissingJSONField(t *testing.T, config map[string]json.RawMessage, field string) {
	t.Helper()
	if _, ok := config[field]; ok {
		t.Fatalf("迁移后不应保留字段 %s", field)
	}
}

func mustMarshalJSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("序列化 JSON 失败: %v", err)
	}
	return encoded
}
