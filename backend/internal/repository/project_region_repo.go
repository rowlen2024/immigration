package repository

import "mygo-immigration/backend/internal/model"

func (r *ProjectRepo) FindPublishedProjectRegions(ids []uint64) ([]model.Project, error) {
	projects := make([]model.Project, 0)
	if len(ids) == 0 {
		return projects, nil
	}
	err := r.db.Select("id", "region_id").Where("status = ? AND id IN ?", 1, ids).
		Preload("Region").Order("id").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepo) ListRegions() ([]model.ProjectRegion, error) {
	var regions []model.ProjectRegion
	err := r.db.Order("sort_order, name, id").Find(&regions).Error
	return regions, err
}

func (r *ProjectRepo) FindRegionByID(id uint64) (*model.ProjectRegion, error) {
	var region model.ProjectRegion
	if err := r.db.Where("id = ?", id).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *ProjectRepo) FindRegionByName(name string) (*model.ProjectRegion, error) {
	var region model.ProjectRegion
	if err := r.db.Where("CAST(name AS BINARY) = CAST(? AS BINARY)", name).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}
