package repository

import (
	"errors"

	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrUnknownProjectCategories = errors.New("选择的项目类目不存在，请刷新后重试")

func (r *ProjectRepo) CreateWithCategories(project *model.Project, categoryIDs []uint64, requestedStatus *int8) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			return err
		}
		// GORM 创建时会用默认值替换显式的 0，在事务内恢复请求状态。
		if requestedStatus != nil && project.Status != *requestedStatus {
			if err := tx.Model(project).UpdateColumn("status", *requestedStatus).Error; err != nil {
				return err
			}
			project.Status = *requestedStatus
		}
		return replaceProjectCategories(tx, project.ID, categoryIDs, nil)
	})
}

func (r *ProjectRepo) UpdateWithCategories(project *model.Project, categoryIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockProjectForCategoryUpdate(tx, project.ID); err != nil {
			return err
		}
		if err := tx.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases", "Advantages", "created_at").Save(project).Error; err != nil {
			return err
		}
		previousIDs, err := findProjectCategoryIDs(tx, project.ID)
		if err != nil {
			return err
		}
		return replaceProjectCategories(tx, project.ID, categoryIDs, previousIDs)
	})
}

func (r *ProjectRepo) deleteWithCategories(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := lockProjectForCategoryUpdate(tx.Unscoped(), id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		categoryIDs, err := findProjectCategoryIDs(tx, id)
		if err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&model.Project{}, id).Error; err != nil {
			return err
		}
		return touchProjectCategories(tx, categoryIDs)
	})
}

func lockProjectForCategoryUpdate(tx *gorm.DB, id uint64) error {
	var project model.Project
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").First(&project, id).Error
}

func findProjectCategoryIDs(tx *gorm.DB, projectID uint64) ([]uint64, error) {
	var ids []uint64
	err := tx.Model(&model.ProjectCategoryMembership{}).Where("project_id = ?", projectID).Pluck("category_id", &ids).Error
	return ids, err
}

func replaceProjectCategories(tx *gorm.DB, projectID uint64, ids, previousIDs []uint64) error {
	// 旧客户端省略 category_ids 时，只更新类目页面的数据时间。
	if ids == nil {
		return touchProjectCategories(tx, previousIDs)
	}
	selected := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		selected[id] = true
	}
	var removed []uint64
	for _, id := range previousIDs {
		if selected[id] {
			delete(selected, id)
		} else {
			removed = append(removed, id)
		}
	}
	membershipsToAdd := make([]model.ProjectCategoryMembership, 0, len(selected))
	newIDs := make([]uint64, 0, len(selected))
	for _, id := range ids {
		if selected[id] {
			newIDs = append(newIDs, id)
			membershipsToAdd = append(membershipsToAdd, model.ProjectCategoryMembership{ProjectID: projectID, CategoryID: id})
			delete(selected, id)
		}
	}
	if len(newIDs) > 0 {
		var found int64
		if err := tx.Model(&model.ProjectCategory{}).Where("id IN ?", newIDs).Count(&found).Error; err != nil {
			return err
		}
		if int64(len(newIDs)) != found {
			return ErrUnknownProjectCategories
		}
		if err := tx.Create(&membershipsToAdd).Error; err != nil {
			return err
		}
	}
	if len(removed) > 0 {
		if err := tx.Where("project_id = ? AND category_id IN ?", projectID, removed).Delete(&model.ProjectCategoryMembership{}).Error; err != nil {
			return err
		}
	}
	return touchProjectCategories(tx, append(previousIDs, newIDs...))
}
