package database

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestProjectCategoryMigration(t *testing.T) {
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
		t.Fatalf("获取测试数据库名称失败: %v", err)
	}
	if !strings.HasSuffix(databaseName, "_migration_test") {
		t.Fatalf("迁移测试只能连接名称以 _migration_test 结尾的数据库，当前为 %q", databaseName)
	}

	run := func(sql string) {
		t.Helper()
		if err := db.Exec(sql).Error; err != nil {
			t.Fatalf("执行 %q 失败: %v", sql, err)
		}
	}
	assertCount := func(sql string, want int64) {
		t.Helper()
		var count int64
		if err := db.Raw(sql).Scan(&count).Error; err != nil || count != want {
			t.Fatalf("执行 %q 得到 %d, %v，期望 %d", sql, count, err, want)
		}
	}
	assertMySQLError := func(sql string, want uint16) {
		t.Helper()
		err := db.Exec(sql).Error
		var mysqlErr *mysqldriver.MySQLError
		if !errors.As(err, &mysqlErr) || mysqlErr.Number != want {
			t.Fatalf("执行 %q 得到 %v，期望 MySQL 错误 %d", sql, err, want)
		}
	}
	reset := func() {
		run("DROP TABLE IF EXISTS navigations, project_category_memberships, project_categories, projects, migrations")
	}
	reset()
	t.Cleanup(reset)
	run("CREATE TABLE projects (id BIGINT UNSIGNED NOT NULL PRIMARY KEY) ENGINE=InnoDB")
	run("CREATE TABLE navigations (id BIGINT UNSIGNED NOT NULL PRIMARY KEY, label VARCHAR(255) NOT NULL) ENGINE=InnoDB")
	run("INSERT INTO projects VALUES (1), (2)")
	run("INSERT INTO navigations VALUES (1, 'existing-link')")

	const migration = "000031_add_project_categories.up.sql"
	migrationsDir := filepath.Join("..", "..", "database", "migrations")
	contents, err := os.ReadFile(filepath.Join(migrationsDir, migration))
	if err != nil {
		t.Fatalf("读取类目迁移文件失败: %v", err)
	}
	testDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(testDir, migration), contents, 0o600); err != nil {
		t.Fatalf("复制类目迁移文件失败: %v", err)
	}
	if err := RunMigrations(db, testDir); err != nil {
		t.Fatalf("执行类目迁移失败: %v", err)
	}
	assertCount("SELECT COUNT(*) FROM project_categories", 0)
	assertCount("SELECT COUNT(*) FROM project_category_memberships", 0)
	assertCount("SELECT COUNT(*) FROM navigations WHERE id = 1 AND label = 'existing-link' AND category_id IS NULL", 1)
	run(`INSERT INTO project_categories (id, name, slug, intro) VALUES
		(11, 'Investment', 'investment', ''), (12, 'Property', 'property', '')`)
	assertCount("SELECT COUNT(*) FROM project_categories WHERE status = 'draft'", 2)
	run(`INSERT INTO project_category_memberships (project_id, category_id, sort_order) VALUES
		(1, 11, 20), (1, 12, 12), (2, 11, 10)`)
	run("UPDATE navigations SET category_id = 11 WHERE id = 1")

	// 模拟 DDL 已提交但执行记录未写入，重跑后必须保留已关联的数据。
	run("DELETE FROM migrations WHERE name = '000031_add_project_categories.up.sql'")
	if err := RunMigrations(db, testDir); err != nil {
		t.Fatalf("重试类目迁移失败: %v", err)
	}
	assertCount("SELECT COUNT(*) FROM project_category_memberships WHERE project_id = 1", 2)
	assertCount("SELECT COUNT(*) FROM project_category_memberships WHERE project_id = 2 AND category_id = 11 AND sort_order = 10", 1)
	assertCount("SELECT COUNT(*) FROM navigations WHERE id = 1 AND category_id = 11", 1)
	assertMySQLError("INSERT INTO project_categories (name, slug, intro) VALUES ('Duplicate', 'investment', '')", 1062)
	assertMySQLError("INSERT INTO project_category_memberships (project_id, category_id) VALUES (1, 11)", 1062)
	assertMySQLError("INSERT INTO project_category_memberships (project_id, category_id) VALUES (999, 11)", 1452)
	assertMySQLError("UPDATE navigations SET category_id = 999 WHERE id = 1", 1452)
	assertMySQLError("DELETE FROM project_categories WHERE id = 11", 1451)

	run("DELETE FROM project_categories WHERE id = 12")
	assertCount("SELECT COUNT(*) FROM project_category_memberships WHERE category_id = 12", 0)
	assertCount("SELECT COUNT(*) FROM projects", 2)
	run("DELETE FROM projects WHERE id = 1")
	assertCount("SELECT COUNT(*) FROM project_category_memberships WHERE project_id = 1", 0)
	assertCount("SELECT COUNT(*) FROM project_categories WHERE id = 11", 1)

	down, err := os.ReadFile(filepath.Join(migrationsDir, "000031_add_project_categories.down.sql"))
	if err != nil {
		t.Fatalf("读取类目回滚文件失败: %v", err)
	}
	for _, statement := range splitStatements(string(down)) {
		run(statement)
	}
	if db.Migrator().HasTable("project_categories") || db.Migrator().HasTable("project_category_memberships") || db.Migrator().HasColumn("navigations", "category_id") {
		t.Fatal("回滚后新增类目表或关联列仍然存在")
	}
	assertCount("SELECT COUNT(*) FROM navigations WHERE id = 1 AND label = 'existing-link'", 1)
	assertCount("SELECT COUNT(*) FROM projects WHERE id = 2", 1)
}
