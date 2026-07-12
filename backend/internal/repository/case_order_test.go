package repository

import (
	"strings"
	"testing"
)

func TestCaseRepoFindAllOrdersPinnedCasesFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &CaseRepo{db: db}

	if _, _, err := repo.FindAll(CaseFilter{Page: 2, PerPage: 10}); err != nil {
		t.Fatalf("生成案例列表 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("案例列表排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
	if strings.Index(sql, expected) > strings.Index(sql, "limit 10") {
		t.Fatalf("案例列表必须先排序再分页，实际 SQL: %s", sql)
	}
}

func TestCaseRepoFindOptionsOrdersPinnedCasesFirst(t *testing.T) {
	db, recorder := newDryRunDB(t)
	repo := &CaseRepo{db: db}

	if _, _, err := repo.FindOptions(CaseFilter{Page: 1, PerPage: 500}); err != nil {
		t.Fatalf("生成案例选项 SQL 失败: %v", err)
	}

	sql := lastSQL(t, recorder)
	expected := "order by is_pinned desc, sort_order asc, id desc"
	if !strings.Contains(sql, expected) {
		t.Fatalf("案例选项排序不正确，期望包含 %q，实际 SQL: %s", expected, sql)
	}
}
