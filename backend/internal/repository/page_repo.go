package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type PageRepo struct {
	db *gorm.DB
}

func (r *PageRepo) FindByID(id uint64) (*model.Page, error) {
	var page model.Page
	err := r.db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindBySlug(slug string) (*model.Page, error) {
	var page model.Page
	err := r.db.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindAll(pageType, search, status string) ([]model.Page, error) {
	var pages []model.Page
	q := r.db.Order("sort_order asc")
	if pageType != "" {
		q = q.Where("page_type = ?", pageType)
	}
	if search != "" {
		q = q.Where("title LIKE ?", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepo) FindAllPublished() ([]model.Page, error) {
	var pages []model.Page
	err := r.db.Where("status = ?", "published").Order("sort_order asc").Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	var page model.Page
	err := r.db.Where("slug = ? AND status = ?", slug, "published").First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindByProjectID(projectID uint64) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Where("project_id = ?", projectID).
		Order("sort_order asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepo) FindAllPaginated(page, perPage int, pageType, search, status string) ([]model.Page, int64, error) {
	var pages []model.Page
	var total int64

	q := r.db.Model(&model.Page{})
	if pageType != "" {
		q = q.Where("page_type = ?", pageType)
	}
	if search != "" {
		q = q.Where("title LIKE ?", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	countQ := q.Session(&gorm.Session{})
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("sort_order asc").Offset(offset).Limit(perPage).Find(&pages).Error
	return pages, total, err
}

func (r *PageRepo) Create(page *model.Page) error {
	return r.db.Create(page).Error
}

func (r *PageRepo) Update(page *model.Page) error {
	return r.db.Omit("created_at").Save(page).Error
}

func (r *PageRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Page{}, id).Error
}

func (r *PageRepo) Count() (int64, error) {
	return CountByModel[model.Page](r.db)
}

func (r *PageRepo) CountByRange(start, end time.Time) (int64, error) {
	return CountByModelRange[model.Page](r.db, start, end)
}

// FindAllCoverImages returns non-empty cover_image values referencing /uploads/ (unscoped).
func (r *PageRepo) FindAllCoverImages() ([]string, error) {
	return PluckUploadsByColumn[model.Page](r.db, "cover_image")
}

// FindAllContents returns content values that contain /uploads/ references (unscoped).
func (r *PageRepo) FindAllContents() ([]string, error) {
	return PluckUploadsByColumn[model.Page](r.db, "content")
}

func (r *PageRepo) Search(keyword string) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("sort_order asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}
