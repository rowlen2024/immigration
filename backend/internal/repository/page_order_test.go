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
