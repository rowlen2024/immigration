package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	projectPinnedUpMigration   = "000026_add_project_pinned.up.sql"
	projectPinnedDownMigration = "000026_add_project_pinned.down.sql"
)

func TestProjectPinnedMigrationFiles(t *testing.T) {
	up := readProjectPinnedMigration(t, projectPinnedUpMigration)
	down := readProjectPinnedMigration(t, projectPinnedDownMigration)

	if !strings.Contains(strings.ToLower(up), "add column `is_pinned` tinyint(1) not null default 0") {
		t.Fatalf("置顶迁移缺少 is_pinned 默认值定义: %s", up)
	}
	if !strings.Contains(strings.ToLower(down), "drop column `is_pinned`") {
		t.Fatalf("置顶回滚迁移缺少 is_pinned 删除语句: %s", down)
	}
}

func TestProjectPinnedMigrationExecution(t *testing.T) {
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

	resetProjectPinnedMigrationTables(t, db)
	t.Cleanup(func() { resetProjectPinnedMigrationTables(t, db) })
	if err := db.Exec("CREATE TABLE projects (id BIGINT UNSIGNED NOT NULL PRIMARY KEY, name VARCHAR(128) NOT NULL) ENGINE=InnoDB").Error; err != nil {
		t.Fatalf("创建项目迁移测试表失败: %v", err)
	}
	if err := db.Exec("INSERT INTO projects (id, name) VALUES (1, 'existing')").Error; err != nil {
		t.Fatalf("插入已有项目失败: %v", err)
	}

	dir := t.TempDir()
	up := readProjectPinnedMigration(t, projectPinnedUpMigration)
	if err := os.WriteFile(filepath.Join(dir, projectPinnedUpMigration), []byte(up), 0o600); err != nil {
		t.Fatalf("准备置顶迁移目录失败: %v", err)
	}
	if err := RunMigrations(db, dir); err != nil {
		t.Fatalf("执行置顶迁移失败: %v", err)
	}

	var isPinned bool
	if err := db.Raw("SELECT is_pinned FROM projects WHERE id = 1").Scan(&isPinned).Error; err != nil {
		t.Fatalf("读取已有项目置顶状态失败: %v", err)
	}
	if isPinned {
		t.Fatal("已有项目迁移后应默认为未置顶")
	}
}

func readProjectPinnedMigration(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", name))
	if err != nil {
		t.Fatalf("读取迁移文件 %s 失败: %v", name, err)
	}
	return string(content)
}

func resetProjectPinnedMigrationTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, table := range []string{"migrations", "projects"} {
		if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			t.Fatalf("清理置顶迁移测试表 %s 失败: %v", table, err)
		}
	}
}
