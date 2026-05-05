package service

import (
	"errors"
	"testing"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// mockFAQRepo implements repository.FAQRepository for testing.
type mockFAQRepo struct {
	findAllFn func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error)
	createFn  func(faq *model.FAQ) error
	updateFn  func(faq *model.FAQ) error
	deleteFn  func(id uint64) error
	searchFn  func(keyword string) ([]model.FAQ, error)
}

func (m *mockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(params)
	}
	return nil, 0, nil
}

func (m *mockFAQRepo) Create(faq *model.FAQ) error {
	if m.createFn != nil {
		return m.createFn(faq)
	}
	return nil
}

func (m *mockFAQRepo) Update(faq *model.FAQ) error {
	if m.updateFn != nil {
		return m.updateFn(faq)
	}
	return nil
}

func (m *mockFAQRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func (m *mockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func TestFAQ_List(t *testing.T) {
	sampleFAQs := []repository.FAQWithProject{
		{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return sampleFAQs, 2, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, err := svc.List(nil, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(faqs) != 2 {
		t.Errorf("expected 2 faqs, got %d", len(faqs))
	}
}

func TestFAQ_List_Error(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.List(nil, nil)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_Create_Success(t *testing.T) {
	created := false
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			created = true
			faq.ID = 100
			return nil
		},
	}

	svc := NewFAQService(repo)

	faq, err := svc.Create(&model.FAQ{Question: "Test Q?", Answer: "Test A."})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !created {
		t.Error("expected Create to be called on repo")
	}
	if faq.ID != 100 {
		t.Errorf("expected ID 100, got %d", faq.ID)
	}
}

func TestFAQ_Create_NilFAQ(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil faq")
	}
}

func TestFAQ_Create_MissingQuestion(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Answer: "Some answer"})
	if err == nil {
		t.Fatal("expected error for missing question")
	}
}

func TestFAQ_Create_MissingAnswer(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Question: "Some question?"})
	if err == nil {
		t.Fatal("expected error for missing answer")
	}
}

func TestFAQ_Update_Success(t *testing.T) {
	updated := false
	repo := &mockFAQRepo{
		updateFn: func(faq *model.FAQ) error {
			updated = true
			return nil
		},
	}

	svc := NewFAQService(repo)

	faq, err := svc.Update(1, &model.FAQ{Question: "Updated Q", Answer: "Updated A"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !updated {
		t.Error("expected Update to be called on repo")
	}
	if faq.ID != 1 {
		t.Errorf("expected ID 1, got %d", faq.ID)
	}
}

func TestFAQ_Update_NilFAQ(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Update(1, nil)
	if err == nil {
		t.Fatal("expected error for nil faq in update")
	}
}

func TestFAQ_Update_ZeroID(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Update(0, &model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestFAQ_Delete_Success(t *testing.T) {
	deleted := false
	repo := &mockFAQRepo{
		deleteFn: func(id uint64) error {
			deleted = true
			return nil
		},
	}

	svc := NewFAQService(repo)

	err := svc.Delete(1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
}

func TestFAQ_Delete_ZeroID(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestFAQ_Delete_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_AdminList_Success(t *testing.T) {
	sampleFAQs := []repository.FAQWithProject{
		{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 3, Question: "Q3", Answer: "A3", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 4, Question: "Q4", Answer: "A4", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			// Simulate pagination: return first page of 2
			return sampleFAQs[:2], 4, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, total, err := svc.AdminList(nil, "", 1, 2)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total 4, got %d", total)
	}
	if len(faqs) != 2 {
		t.Errorf("expected 2 faqs on page 1 with perPage=2, got %d", len(faqs))
	}
}

func TestFAQ_AdminList_Page2(t *testing.T) {
	sampleFAQs := []repository.FAQWithProject{
		{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 3, Question: "Q3", Answer: "A3", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			// Simulate pagination: return page 2 (only 1 item remaining)
			return sampleFAQs[2:], 3, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, total, err := svc.AdminList(nil, "", 2, 2)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(faqs) != 1 {
		t.Errorf("expected 1 faq on page 2 with perPage=2, got %d", len(faqs))
	}
}

func TestFAQ_AdminList_BeyondRange(t *testing.T) {
	sampleFAQs := []repository.FAQWithProject{
		{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return sampleFAQs, 1, nil
		},
	}

	svc := NewFAQService(repo)

	// Note: Pagination now happens in the repo, so beyond-range is handled at DB level
	_, total, err := svc.AdminList(nil, "", 10, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
}

func TestFAQ_AdminList_DefaultPagination(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, nil
		},
	}

	svc := NewFAQService(repo)
	_, _, err := svc.AdminList(nil, "", 0, 0)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_AdminList_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, _, err := svc.AdminList(nil, "", 1, 10)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_Create_SpecialCharacters(t *testing.T) {
	var savedFAQ *model.FAQ
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			savedFAQ = faq
			faq.ID = 1
			return nil
		},
	}

	svc := NewFAQService(repo)

	question := "What about <script>alert('xss')</script> & special chars?"
	answer := "Answer with <b>bold</b> and &amp; encoding"
	faq, err := svc.Create(&model.FAQ{Question: question, Answer: answer})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if savedFAQ == nil {
		t.Fatal("expected Create to be called on repo")
	}
	if savedFAQ.Question != question {
		t.Errorf("expected question to be saved as-is")
	}
	if faq.ID != 1 {
		t.Errorf("expected ID 1, got %d", faq.ID)
	}
}

func TestFAQ_Create_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_Update_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		updateFn: func(faq *model.FAQ) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.Update(1, &model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_List_Empty(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, err := svc.List(nil, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(faqs) != 0 {
		t.Errorf("expected 0 faqs, got %d", len(faqs))
	}
}

func TestFAQ_List_WithProjectData(t *testing.T) {
	pid := uint64(5)
	sampleFAQs := []repository.FAQWithProject{
		{
			FAQ:         model.FAQ{ID: 1, Question: "Q1", Answer: "A1", ProjectID: &pid, IsGlobal: false, SortOrder: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			ProjectName: "EB-5 Visa",
			ProjectSlug: "eb5-visa",
		},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return sampleFAQs, 1, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, err := svc.List(&pid, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(faqs) != 1 {
		t.Fatalf("expected 1 faq, got %d", len(faqs))
	}
	resp := faqs[0]
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
	if resp.ProjectName != "EB-5 Visa" {
		t.Errorf("expected project name 'EB-5 Visa', got '%s'", resp.ProjectName)
	}
	if resp.ProjectSlug != "eb5-visa" {
		t.Errorf("expected project slug 'eb5-visa', got '%s'", resp.ProjectSlug)
	}
	if resp.Question != "Q1" {
		t.Errorf("expected question 'Q1', got '%s'", resp.Question)
	}

	// Verify it's a dto.FAQResponse (with string timestamps)
	var _ dto.FAQResponse = resp
	if resp.CreatedAt == "" {
		t.Error("expected non-empty created_at string")
	}
}
