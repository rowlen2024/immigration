package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Project       *ProjectRepo
	User          *UserRepo
	FAQ           *FAQRepo
	Case          *CaseRepo
	Page          *PageRepo
	Lead          *LeadRepo
	HomeConfig    *HomeConfigRepo
	Media         *MediaRepo
	Nav           *NavRepo
	Requirement   *RequirementRepo
	CostItem      *CostItemRepo
	TimelinePhase  *TimelinePhaseRepo
	CompareConfig   *CompareConfigRepo
	ProjectAdvantage *ProjectAdvantageRepo
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		Project:       &ProjectRepo{db: db},
		User:          &UserRepo{db: db},
		FAQ:           &FAQRepo{db: db},
		Case:          &CaseRepo{db: db},
		Page:          &PageRepo{db: db},
		Lead:          &LeadRepo{db: db},
		HomeConfig:    &HomeConfigRepo{db: db},
		Media:         &MediaRepo{db: db},
		Nav:           &NavRepo{db: db},
		Requirement:   &RequirementRepo{db: db},
		CostItem:      &CostItemRepo{db: db},
		TimelinePhase:  &TimelinePhaseRepo{db: db},
		CompareConfig:  &CompareConfigRepo{db: db},
			ProjectAdvantage: &ProjectAdvantageRepo{db: db},
	}
}
