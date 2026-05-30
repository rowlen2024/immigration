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
	findByIDFn             func(id uint64) (*model.FAQ, error)
	findAllFn              func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error)
	findDistinctProjectsFn func() ([]model.Project, error)
	createFn               func(faq *model.FAQ) error
	updateFn               func(faq *model.FAQ) error
	deleteFn               func(id uint64) error
	deleteByProjectIDFn    func(projectID uint64) error
	searchFn               func(keyword string) ([]model.FAQ, error)
}

func (m *mockFAQRepo) FindByID(id uint64) (*model.FAQ, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *mockFAQRepo) FindDistinctProjects() ([]model.Project, error) {
	if m.findDistinctProjectsFn != nil {
		return m.findDistinctProjectsFn()
	}
	return nil, nil
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

func (m *mockFAQRepo) DeleteByProjectID(projectID uint64) error {
	if m.deleteByProjectIDFn != nil {
		return m.deleteByProjectIDFn(projectID)
	}
	return nil
}

func (m *mockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func TestFAQ_List_Success(t *testing.T) {
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

	faqs, total, err := svc.List(dto.FAQListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
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

	_, _, err := svc.List(dto.FAQListRequest{})
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

	faqs, total, err := svc.List(dto.FAQListRequest{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
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

	faqs, _, err := svc.List(dto.FAQListRequest{ProjectID: &pid})
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
}

func TestFAQ_List_WithIsGlobal(t *testing.T) {
	isGlobal := true
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.IsGlobal == nil || *params.IsGlobal != true {
				t.Error("expected IsGlobal to be true")
			}
			return nil, 0, nil
		},
	}

	svc := NewFAQService(repo)
	_, _, err := svc.List(dto.FAQListRequest{IsGlobal: &isGlobal})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_List_Paginated(t *testing.T) {
	sampleFAQs := []repository.FAQWithProject{
		{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return sampleFAQs, 4, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, total, err := svc.List(dto.FAQListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 2},
	})
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

func TestFAQ_Update_Success(t *testing.T) {
	updated := false
	repo := &mockFAQRepo{
		findByIDFn: func(id uint64) (*model.FAQ, error) {
			return &model.FAQ{ID: id, Question: "Old Q", Answer: "Old A"}, nil
		},
		updateFn: func(faq *model.FAQ) error {
			updated = true
			return nil
		},
	}

	svc := NewFAQService(repo)

	faq, err := svc.Update(1, dto.UpdateFAQRequest{Question: "Updated Q", Answer: "Updated A"})
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

	_, err := svc.Update(1, dto.UpdateFAQRequest{})
	if err == nil {
		t.Fatal("expected error for nil faq in update")
	}
}

func TestFAQ_Update_ZeroID(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Update(0, dto.UpdateFAQRequest{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestFAQ_Update_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		updateFn: func(faq *model.FAQ) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.Update(1, dto.UpdateFAQRequest{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error from repo")
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
