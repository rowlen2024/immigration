package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestProjectRegionMigration(t *testing.T) {
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
			t.Fatalf("执行迁移测试语句失败: %v", err)
		}
	}
	reset := func() {
		run("DROP TABLE IF EXISTS projects")
		run("DROP TABLE IF EXISTS project_regions")
		run("DROP TABLE IF EXISTS migrations")
	}
	reset()
	t.Cleanup(reset)
	run(`CREATE TABLE projects (
		id BIGINT UNSIGNED NOT NULL PRIMARY KEY,
		slug VARCHAR(64) NOT NULL,
		country VARCHAR(64) NOT NULL,
		deleted_at DATETIME NULL,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	run(`INSERT INTO projects (id, slug, country, deleted_at, updated_at) VALUES
		(1, 'cies', '中国香港', NULL, '2026-01-01 00:00:00'),
		(2, 'hong-kong-other', ' 中国（香港特别行政区） ', NULL, '2026-01-01 00:00:00'),
		(3, 'xiang-gang-gao-duan-ren-cai', '中国', NULL, '2026-01-01 00:00:00'),
		(4, 'china-other', '中国', NULL, '2026-01-01 00:00:00'),
		(5, 'unknown', '待确认地区', NULL, '2026-01-01 00:00:00'),
		(6, 'deleted', '美国', CURRENT_TIMESTAMP, '2026-01-01 00:00:00')`)

	const migration = "000030_add_project_regions.up.sql"
	contents, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", migration))
	if err != nil {
		t.Fatalf("读取地区迁移文件失败: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, migration), contents, 0o600); err != nil {
		t.Fatalf("复制地区迁移文件失败: %v", err)
	}
	if err := RunMigrations(db, dir); err != nil {
		t.Fatalf("执行地区迁移失败: %v", err)
	}

	var chinaCode string
	if err := db.Raw("SELECT r.code FROM projects p JOIN project_regions r ON r.id = p.region_id WHERE p.id = 4").Scan(&chinaCode).Error; err != nil || chinaCode != "CN" {
		t.Fatalf("中国地区回填结果错误: %q, %v", chinaCode, err)
	}

	// 模拟 DDL 已提交但迁移记录未写入，重试不应修改已回填的地区。
	run(`UPDATE projects SET region_id = (SELECT id FROM project_regions WHERE code = 'US'),
		updated_at = updated_at WHERE id = 4`)
	run("DELETE FROM migrations WHERE name = '000030_add_project_regions.up.sql'")
	if err := RunMigrations(db, dir); err != nil {
		t.Fatalf("重试地区迁移失败: %v", err)
	}

	var rows []struct {
		RegionCode string
		Country    string
		Updated    string
	}
	if err := db.Raw(`SELECT COALESCE(r.code, '') AS region_code, p.country,
		DATE_FORMAT(p.updated_at, '%Y-%m-%d %H:%i:%s') AS updated
		FROM projects p LEFT JOIN project_regions r ON r.id = p.region_id ORDER BY p.id`).Scan(&rows).Error; err != nil {
		t.Fatalf("读取已迁移项目失败: %v", err)
	}
	codes := []string{"HK", "HK", "", "US", "", ""}
	countries := []string{"中国香港", " 中国（香港特别行政区） ", "中国", "中国", "待确认地区", "美国"}
	if len(rows) != len(codes) {
		t.Fatalf("项目数量错误: %d", len(rows))
	}
	for i, row := range rows {
		if row.RegionCode != codes[i] || row.Country != countries[i] || row.Updated != "2026-01-01 00:00:00" {
			t.Errorf("第 %d 个项目回填错误: %+v, 期望地区 %q、文案 %q 并保留时间", i+1, row, codes[i], countries[i])
		}
	}

	down, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", "000030_add_project_regions.down.sql"))
	if err != nil {
		t.Fatalf("读取地区回滚文件失败: %v", err)
	}
	for _, statement := range splitStatements(string(down)) {
		run(statement)
	}
	if db.Migrator().HasTable("project_regions") || db.Migrator().HasColumn("projects", "region_id") {
		t.Fatal("回滚后地区表或关联列仍然存在")
	}
}
