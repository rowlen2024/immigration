package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	db               *gorm.DB
	Project          *ProjectRepo
	User             *UserRepo
	RBAC             *RBACRepo
	FAQ              *FAQRepo
	Case             *CaseRepo
	Page             *PageRepo
	Lead             *LeadRepo
	Lawyer           *LawyerRepo
	HomeConfig       *HomeConfigRepo
	Media            *MediaRepo
	Nav              *NavRepo
	Requirement      *RequirementRepo
	CostItem         *CostItemRepo
	TimelinePhase    *TimelinePhaseRepo
	CompareConfig    *CompareConfigRepo
	ProjectAdvantage *ProjectAdvantageRepo
	Milestone        *MilestoneRepo
	Testimonial      *TestimonialRepo
	PublicVersion    *PublicVersionRepo
}

func New(db *gorm.DB) *Repository {
	SetDefaultDB(db)
	return &Repository{
		db:               db,
		Project:          &ProjectRepo{db: db},
		User:             &UserRepo{db: db},
		RBAC:             &RBACRepo{db: db},
		FAQ:              &FAQRepo{db: db},
		Case:             &CaseRepo{db: db},
		Page:             &PageRepo{db: db},
		Lead:             &LeadRepo{db: db},
		Lawyer:           &LawyerRepo{db: db},
		HomeConfig:       &HomeConfigRepo{db: db},
		Media:            &MediaRepo{db: db},
		Nav:              &NavRepo{db: db},
		Requirement:      &RequirementRepo{db: db},
		CostItem:         &CostItemRepo{db: db},
		TimelinePhase:    &TimelinePhaseRepo{db: db},
		CompareConfig:    &CompareConfigRepo{db: db},
		ProjectAdvantage: &ProjectAdvantageRepo{db: db},
		Milestone:        &MilestoneRepo{db: db},
		Testimonial:      &TestimonialRepo{db: db},
		PublicVersion:    &PublicVersionRepo{db: db},
	}
}

// Tx executes fn within a database transaction, passing a new Repository whose
// repos all share the same tx. If fn returns an error the tx is rolled back.
func (r *Repository) Tx(fn func(txRepo *Repository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	txRepo := New(tx)
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
