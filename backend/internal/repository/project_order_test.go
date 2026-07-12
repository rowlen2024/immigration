package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type sqlRecorder struct {
	statements []string
}

func (r *sqlRecorder) LogMode(logger.LogLevel) logger.Interface      { return r }
func (r *sqlRecorder) Info(context.Context, string, ...interface{})  {}
func (r *sqlRecorder) Warn(context.Context, string, ...interface{})  {}
func (r *sqlRecorder) Error(context.Context, string, ...interface{}) {}
func (r *sqlRecorder) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	sql, _ := fc()
	r.statements = append(r.statements, sql)
}

func newDryRunDB(t *testing.T) (*gorm.DB, *sqlRecorder) {
	t.Helper()
	recorder := &sqlRecorder{}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               recorder,
	})
	if err != nil {
		t.Fatalf("创建 DryRun 数据库失败: %v", err)
	}
	return db, recorder
}

func normalizeSQL(sql string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.ReplaceAll(sql, "`", "")), " "))
}

func lastSQL(t *testing.T, recorder *sqlRecorder) string {
	t.Helper()
	if len(recorder.statements) == 0 {
		t.Fatal("未捕获到 SQL")
	}
	return normalizeSQL(recorder.statements[len(recorder.statements)-1])
}

func TestProjectRepoFindAllOrdersPinnedProjectsFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &ProjectRepo{db: db}

	if _, _, err := repo.FindAll(ProjectFilter{Page: 2, PerPage: 10}); err != nil {
		t.Fatalf("生成项目列表 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("项目列表排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
	if strings.Index(sql, expected) > strings.Index(sql, "limit 10") {
		t.Fatalf("项目列表必须先排序再分页，实际 SQL: %s", sql)
	}
}

func TestProjectRepoFindAllFiltersByCountry(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &ProjectRepo{db: db}

	if _, _, err := repo.FindAll(ProjectFilter{Country: "Canada"}); err != nil {
		t.Fatalf("failed to generate project list SQL: %v", err)
	}

	if len(recorder.statements) != 2 {
		t.Fatalf("expected count and find SQL, got %d statements", len(recorder.statements))
	}
	for _, statement := range recorder.statements {
		sql := normalizeSQL(statement)
		if !strings.Contains(sql, "where country like '%canada%'") {
			t.Fatalf("expected country filter in count and find SQL, got: %s", sql)
		}
	}
}

func TestProjectRepoFindAllCombinesCountryAndNameFilters(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &ProjectRepo{db: db}

	if _, _, err := repo.FindAll(ProjectFilter{Country: "Canada", Name: "Invest", Page: 1, PerPage: 12}); err != nil {
		t.Fatalf("failed to generate project list SQL: %v", err)
	}

	if len(recorder.statements) != 2 {
		t.Fatalf("expected count and find SQL, got %d statements", len(recorder.statements))
	}
	for _, statement := range recorder.statements {
		sql := normalizeSQL(statement)
		if !strings.Contains(sql, "name like '%invest%'") || !strings.Contains(sql, "country like '%canada%'") {
			t.Fatalf("expected name and country filters in count and find SQL, got: %s", sql)
		}
	}
	findSQL := normalizeSQL(recorder.statements[len(recorder.statements)-1])
	if !strings.Contains(findSQL, "order by is_pinned desc, sort_order asc, id desc") {
		t.Fatalf("expected existing pinned ordering, got: %s", findSQL)
	}
}

func TestProjectRepoFindOptionsOrdersPinnedProjectsFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &ProjectRepo{db: db}

	if _, _, err := repo.FindOptions(ProjectFilter{Page: 1, PerPage: 500}); err != nil {
		t.Fatalf("生成项目选项 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("项目选项排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
}

func TestFAQRepoFindDistinctProjectsOrdersPinnedProjectsFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &FAQRepo{db: db}

	if _, err := repo.FindDistinctProjects(); err != nil {
		t.Fatalf("生成 FAQ 项目筛选 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by projects.is_pinned desc, projects.sort_order asc, projects.id asc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("FAQ 项目筛选排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
}
