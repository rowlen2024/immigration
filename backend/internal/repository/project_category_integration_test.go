package repository

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/database"
	"mygo-immigration/backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestProjectCategoryMembershipsMySQL(t *testing.T) {
	dsn := os.Getenv("PROJECT_CATEGORY_TEST_DSN")
	if dsn == "" {
		t.Skip("未设置 PROJECT_CATEGORY_TEST_DSN，跳过类目写入集成测试")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	var schema string
	if err := db.Raw("SELECT DATABASE()").Scan(&schema).Error; err != nil || !strings.HasSuffix(schema, "_migration_test") {
		t.Fatalf("类目写入测试仅支持隔离数据库，实际为 %q，错误 %v", schema, err)
	}
	if err := database.RunMigrations(db, filepath.Join("..", "..", "database", "migrations")); err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Cleanup(func() { tx.Rollback() })
	projects := &ProjectRepo{db: tx}
	categories := &ProjectCategoryRepo{db: tx}

	public := model.ProjectCategory{Name: "分类测试", Slug: "category-test-public", Status: model.ProjectCategoryPublished, SEOTitle: "搜索标题"}
	draft := model.ProjectCategory{Name: "分类测试草稿", Slug: "category-test-draft", Status: model.ProjectCategoryDraft}
	empty := model.ProjectCategory{Name: "分类测试空类目", Slug: "category-test-empty", Status: model.ProjectCategoryPublished}
	for _, category := range []*model.ProjectCategory{&public, &draft, &empty} {
		if err := categories.Create(category); err != nil {
			t.Fatal(err)
		}
	}
	first := model.Project{Name: "First project", Slug: "category-test-first"}
	second := model.Project{Name: "Second project", Slug: "category-test-second", Status: 1}
	if err := projects.CreateWithCategories(&first, []uint64{public.ID, draft.ID, public.ID}, nil); err != nil {
		t.Fatal(err)
	}
	if first.Status != 1 {
		t.Fatalf("legacy default status = %d, want published", first.Status)
	}
	if err := projects.CreateWithCategories(&second, []uint64{public.ID}, nil); err != nil {
		t.Fatal(err)
	}
	zero := int8(0)
	privateProject := model.Project{Name: "Private project", Slug: "category-test-private"}
	if err := projects.CreateWithCategories(&privateProject, []uint64{public.ID}, &zero); err != nil {
		t.Fatal(err)
	}
	stored, err := projects.FindByID(privateProject.ID)
	if err != nil || privateProject.Status != 0 || stored.Status != 0 {
		t.Fatalf("explicit unpublished status not preserved: %+v, %v", stored, err)
	}
	unassigned := model.Project{Name: "Unassigned project", Slug: "category-test-unassigned"}
	if err := projects.CreateWithCategories(&unassigned, nil, &zero); err != nil {
		t.Fatal(err)
	}
	stored, err = projects.FindByID(unassigned.ID)
	if err != nil || stored.Status != 0 {
		t.Fatalf("explicit unpublished status lost without memberships: %+v, %v", stored, err)
	}
	adminRows, total, err := categories.FindAll(ProjectCategoryFilter{Name: "分类测试", Page: 1, PerPage: 10})
	if err != nil || total != 3 || len(adminRows) != 3 || adminRows[0].ID != public.ID || adminRows[0].ProjectCount != 3 {
		t.Fatalf("admin list did not count all assignments: %+v, %d, %v", adminRows, total, err)
	}
	if err := categories.UpdateProjectOrder(public.ID, second.ID, -10); err != nil {
		t.Fatal(err)
	}
	if err := categories.UpdateProjectOrder(public.ID, first.ID, 20); err != nil {
		t.Fatal(err)
	}
	rows, err := categories.FindPublishedProjects(public.ID)
	if err != nil || len(rows) != 2 || rows[0].ID != second.ID || rows[1].ID != first.ID {
		t.Fatalf("public projects in unexpected order: %+v, %v", rows, err)
	}
	options, err := categories.FindPublicSummaries()
	if err != nil || len(options) != 1 || options[0].ID != public.ID || options[0].ProjectCount != 2 || options[0].UpdatedAt.IsZero() {
		t.Fatalf("public categories include a draft or empty item: %+v, %v", options, err)
	}
	if _, err := categories.FindPublishedBySlug(draft.Slug); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("draft was publicly readable: %v", err)
	}
	if _, err := categories.FindPublishedBySlug(strings.ToUpper(public.Slug)); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("non-canonical slug was publicly readable: %v", err)
	}
	if dbRow, err := categories.FindPublishedBySlug(public.Slug); err != nil || dbRow.SEOTitle != public.SEOTitle {
		t.Fatalf("public SEO data missing: %+v, %v", dbRow, err)
	}

	first.Name = "Edited by old client"
	if err := projects.Update(&first); err != nil {
		t.Fatal(err)
	}
	assignments, err := categories.AssignmentsByProjectIDs([]uint64{first.ID, second.ID})
	if err != nil || len(assignments) != 3 {
		t.Fatalf("old update did not preserve memberships: %+v, %v", assignments, err)
	}
	if err := projects.UpdateWithCategories(&first, []uint64{public.ID}); err != nil {
		t.Fatal(err)
	}
	members, total, err := categories.FindProjects(public.ID, 1, 10)
	if err != nil || total != 3 || len(members) != 3 || members[2].SortOrder != 20 {
		t.Fatalf("retained membership lost its order: %+v, %d, %v", members, total, err)
	}

	first.Name = "Failed edit"
	if err := projects.UpdateWithCategories(&first, []uint64{empty.ID + 100}); !errors.Is(err, ErrUnknownProjectCategories) {
		t.Fatalf("unknown category did not reject project edit: %v", err)
	}
	if item, err := projects.FindByID(first.ID); err != nil || item.Name != "Edited by old client" {
		t.Fatalf("failed category update saved project changes: %+v, %v", item, err)
	}
	if associations, err := categories.AssignmentsByProjectIDs([]uint64{first.ID}); err != nil || len(associations) != 1 || associations[0].ID != public.ID {
		t.Fatalf("failed category update changed existing membership: %+v, %v", associations, err)
	}
	if err := projects.UpdateWithCategories(&first, []uint64{}); err != nil {
		t.Fatal(err)
	}
	if associations, err := categories.AssignmentsByProjectIDs([]uint64{first.ID}); err != nil || len(associations) != 0 {
		t.Fatalf("explicit clear retained membership: %+v, %v", associations, err)
	}

	second.Status = 0
	if err := projects.Update(&second); err != nil {
		t.Fatal(err)
	}
	if rows, err := categories.FindPublishedProjects(public.ID); err != nil || len(rows) != 0 {
		t.Fatalf("private project was exposed: %+v, %v", rows, err)
	}
	if options, err := categories.FindPublicSummaries(); err != nil || len(options) != 0 {
		t.Fatalf("category with no visible projects stayed in public list: %+v, %v", options, err)
	}
}
