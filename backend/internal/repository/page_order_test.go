package repository

import (
	"strings"
	"testing"
)

func TestPageRepoFindAllOrdersPinnedPagesFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &PageRepo{db: db}

	if _, _, err := repo.FindAll(PageFilter{Page: 2, PerPage: 10}); err != nil {
		t.Fatalf("生成页面列表 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("页面列表排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
	if strings.Index(sql, expected) > strings.Index(sql, "limit 10") {
		t.Fatalf("页面列表必须先排序再分页，实际 SQL: %s", sql)
	}
}

func TestPageRepoFindOptionsOrdersPinnedPagesFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &PageRepo{db: db}

	if _, _, err := repo.FindOptions(PageFilter{Page: 1, PerPage: 500}); err != nil {
		t.Fatalf("生成页面选项 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("页面选项排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
}

func TestPageRepoFindRelatedBySlugBuildsScopedDistinctQuery(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &PageRepo{db: db}

	if _, err := repo.FindRelatedBySlug("current-news", 4); err != nil {
		t.Fatalf("failed to build related pages SQL: %v", err)
	}

	var sql string
	for _, statement := range recorder.statements {
		normalized := normalizeSQL(statement)
		if strings.Contains(normalized, "select distinct pages.*") {
			sql = normalized
			break
		}
	}
	if sql == "" {
		t.Fatal("related project query was not generated")
	}
	checks := []string{
		"select distinct pages.*",
		"join project_news on project_news.page_id = pages.id",
		"project_news.project_id in (select project_id from project_news where page_id = 0)",
		"pages.id <> 0",
		"pages.status = 'published'",
		"pages.page_type = 'news'",
		"pages.deleted_at is null",
		"order by pages.is_pinned desc, pages.created_at desc, pages.id desc",
		"limit 4",
	}
	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q, got: %s", check, sql)
		}
	}
}
