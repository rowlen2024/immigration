package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/database"
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

func TestPageRepoFindRelatedBySlugMySQL(t *testing.T) {
	db, err := database.InitMySQL(repositoryTestConfig(t))
	if err != nil {
		t.Fatalf("connect to MySQL: %v", err)
	}
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() { tx.Rollback() })

	prefix := fmt.Sprintf("related_test_%d", time.Now().UnixNano())
	projects := []model.Project{
		{Slug: prefix + "_p1", Name: "Related Test P1"},
		{Slug: prefix + "_p2", Name: "Related Test P2"},
		{Slug: prefix + "_unrelated", Name: "Related Test Unrelated"},
	}
	if err := tx.Create(&projects).Error; err != nil {
		t.Fatalf("create projects: %v", err)
	}

	base := time.Date(2099, 7, 12, 0, 0, 0, 0, time.UTC)
	pages := []model.Page{
		{Title: "Current", Slug: prefix + "_current", Status: "published", PageType: "news", IsPinned: true, CreatedAt: base.Add(10 * time.Hour)},
		{Title: "Overlap", Slug: prefix + "_overlap", Status: "published", PageType: "news", CreatedAt: base.Add(5 * time.Hour)},
		{Title: "Pinned", Slug: prefix + "_pinned", Status: "published", PageType: "news", IsPinned: true, CreatedAt: base.Add(time.Hour)},
		{Title: "Newest", Slug: prefix + "_newest", Status: "published", PageType: "news", CreatedAt: base.Add(5 * time.Hour)},
		{Title: "Middle", Slug: prefix + "_middle", Status: "published", PageType: "news", CreatedAt: base.Add(2 * time.Hour)},
		{Title: "Old", Slug: prefix + "_old", Status: "published", PageType: "news", CreatedAt: base.Add(-time.Hour)},
		{Title: "Draft", Slug: prefix + "_draft", Status: "draft", PageType: "news", Tags: []string{"visa", "investment"}, IsPinned: true, CreatedAt: base.Add(30 * time.Hour)},
		{Title: "Not News", Slug: prefix + "_page", Status: "published", PageType: "default", Tags: []string{"visa", "investment"}, IsPinned: true, CreatedAt: base.Add(29 * time.Hour)},
		{Title: "Deleted", Slug: prefix + "_deleted", Status: "published", PageType: "news", Tags: []string{"visa", "investment"}, IsPinned: true, CreatedAt: base.Add(28 * time.Hour)},
		{Title: "Unrelated", Slug: prefix + "_unrelated", Status: "published", PageType: "news", CreatedAt: base.Add(9 * time.Hour)},
		{Title: "No Projects", Slug: prefix + "_no_projects", Status: "published", PageType: "news", CreatedAt: base},
		{Title: "Tag Current", Slug: prefix + "_tag_current", Status: "published", PageType: "news", Tags: []string{"visa", "investment"}, CreatedAt: base.Add(20 * time.Hour)},
		{Title: "Tag High", Slug: prefix + "_tag_high", Status: "published", PageType: "news", Tags: []string{"visa", "investment"}, CreatedAt: base.Add(time.Hour)},
		{Title: "Tag Low", Slug: prefix + "_tag_low", Status: "published", PageType: "news", Tags: []string{"visa"}, IsPinned: true, CreatedAt: base.Add(19 * time.Hour)},
		{Title: "Tag None", Slug: prefix + "_tag_none", Status: "published", PageType: "news", Tags: []string{"visa"}, IsPinned: true, CreatedAt: base.Add(18 * time.Hour)},
	}
	if err := tx.Create(&pages).Error; err != nil {
		t.Fatalf("create pages: %v", err)
	}

	links := []model.ProjectNews{
		{ProjectID: projects[0].ID, PageID: pages[0].ID},
		{ProjectID: projects[1].ID, PageID: pages[0].ID},
		{ProjectID: projects[0].ID, PageID: pages[1].ID},
		{ProjectID: projects[1].ID, PageID: pages[1].ID},
		{ProjectID: projects[0].ID, PageID: pages[2].ID},
		{ProjectID: projects[1].ID, PageID: pages[3].ID},
		{ProjectID: projects[0].ID, PageID: pages[4].ID},
		{ProjectID: projects[0].ID, PageID: pages[5].ID},
		{ProjectID: projects[0].ID, PageID: pages[6].ID},
		{ProjectID: projects[0].ID, PageID: pages[7].ID},
		{ProjectID: projects[0].ID, PageID: pages[8].ID},
		{ProjectID: projects[2].ID, PageID: pages[11].ID},
		{ProjectID: projects[2].ID, PageID: pages[14].ID},
	}
	if err := tx.Create(&links).Error; err != nil {
		t.Fatalf("create project_news: %v", err)
	}
	if err := tx.Delete(&pages[8]).Error; err != nil {
		t.Fatalf("soft delete page: %v", err)
	}

	repo := &PageRepo{db: tx}
	got, err := repo.FindRelatedBySlug(pages[0].Slug, 4)
	if err != nil {
		t.Fatalf("find related pages: %v", err)
	}
	wantIDs := []uint64{pages[2].ID, pages[3].ID, pages[1].ID, pages[4].ID}
	if len(got) != len(wantIDs) {
		t.Fatalf("expected %d related pages, got %d: %#v", len(wantIDs), len(got), got)
	}
	for i, wantID := range wantIDs {
		if got[i].ID != wantID {
			t.Fatalf("unexpected result at %d: want id %d, got id %d", i, wantID, got[i].ID)
		}
	}
	if got[1].CreatedAt != got[2].CreatedAt || got[1].ID <= got[2].ID {
		t.Fatalf("same-created-at pages must be ordered by id desc, got %#v then %#v", got[1], got[2])
	}

	fallback, err := repo.FindRelatedBySlug(pages[10].Slug, 4)
	if err != nil {
		t.Fatalf("find fallback pages without projects or tags: %v", err)
	}
	wantFallback := []uint64{pages[13].ID, pages[14].ID, pages[0].ID, pages[2].ID}
	if len(fallback) != len(wantFallback) {
		t.Fatalf("expected %d fallback pages, got %#v", len(wantFallback), fallback)
	}
	for i, wantID := range wantFallback {
		if fallback[i].ID != wantID {
			t.Fatalf("unexpected fallback at %d: want %d, got %d", i, wantID, fallback[i].ID)
		}
	}

	tagged, err := repo.FindRelatedBySlug(pages[11].Slug, 4)
	if err != nil {
		t.Fatalf("find tag-related pages: %v", err)
	}
	wantTagged := []uint64{pages[14].ID, pages[12].ID, pages[13].ID, pages[0].ID}
	if len(tagged) != len(wantTagged) {
		t.Fatalf("expected tag and fallback pages, got %#v", tagged)
	}
	for i, wantID := range wantTagged {
		if tagged[i].ID != wantID {
			t.Fatalf("unexpected tag fallback at %d: want %d, got %d", i, wantID, tagged[i].ID)
		}
	}

	_, err = repo.FindRelatedBySlug(prefix+"_missing", 4)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}

	projectSummaries, err := repo.FindProjectsByPageID(pages[0].ID)
	if err != nil {
		t.Fatalf("find page projects: %v", err)
	}
	if len(projectSummaries) != 2 || projectSummaries[0].ID != projects[0].ID || projectSummaries[1].ID != projects[1].ID {
		t.Fatalf("unexpected project summaries: %#v", projectSummaries)
	}
}

func repositoryTestConfig(t *testing.T) *config.Config {
	t.Helper()
	file, err := os.Open("../../../.env")
	if err != nil {
		t.Fatalf("open project .env: %v", err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(strings.TrimSpace(scanner.Text()), "=", 2)
		if len(parts) == 2 {
			values[parts[0]] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read project .env: %v", err)
	}
	return &config.Config{
		DBHost: values["DB_HOST"], DBPort: values["DB_PORT"], DBUser: values["DB_USER"],
		DBPassword: values["DB_PASSWORD"], DBName: values["DB_NAME"],
	}
}
