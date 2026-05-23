package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type LawyerRepo struct {
	db *gorm.DB
}

func (r *LawyerRepo) FindAll() ([]model.Lawyer, error) {
	var items []model.Lawyer
	err := r.db.Order("sort_order ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *LawyerRepo) FindPaginated(page, perPage int, search string) ([]model.Lawyer, int64, error) {
	var items []model.Lawyer
	var total int64

	q := r.db.Model(&model.Lawyer{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("sort_order ASC, id ASC").Offset(offset).Limit(perPage).Find(&items).Error
	return items, total, err
}

func (r *LawyerRepo) FindByID(id uint64) (*model.Lawyer, error) {
	var item model.Lawyer
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *LawyerRepo) Create(item *model.Lawyer) error {
	return r.db.Create(item).Error
}

func (r *LawyerRepo) Update(item *model.Lawyer) error {
	return r.db.Omit("created_at").Save(item).Error
}

func (r *LawyerRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Lawyer{}, id).Error
}

// FindAllPhotoURLs returns non-empty photo_url values referencing /uploads/ (unscoped).
func (r *LawyerRepo) FindAllPhotoURLs() ([]string, error) {
	var urls []string
	err := r.db.Unscoped().Model(&model.Lawyer{}).
		Where("photo_url LIKE ?", "%/uploads/%").
		Pluck("photo_url", &urls).Error
	return urls, err
}
