package service

import (
	"reflect"
	"strconv"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/repository"
)

type Service struct {
	repo          *repository.Repository
	Project       *ProjectService
	Auth          *AuthService
	User          *UserService
	FAQ           *FAQService
	Page          *PageService
	Case          *CaseService
	Lead          *LeadService
	Lawyer        *LawyerService
	HomeConfig    *HomeConfigService
	Media         *MediaService
	Nav           *NavService
	Search        *SearchService
	Requirement   *RequirementService
	CostItem      *CostItemService
	TimelinePhase *TimelinePhaseService
	CompareConfig *CompareConfigService
	Advantage     *ProjectAdvantageService
	Testimonial   *TestimonialService
	PublicVersion *PublicVersionService
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	svc := &Service{
		repo:       repo,
		Project:    NewProjectService(repo.Project, repo.Nav),
		Auth:       NewAuthService(repo.User, cfg),
		User:       NewUserService(repo.User),
		FAQ:        NewFAQService(repo.FAQ),
		Page:       NewPageService(repo.Page, repo.Nav),
		Case:       NewCaseService(repo.Case, nil),
		Lead:       NewLeadService(repo.Lead),
		Lawyer:     NewLawyerService(repo.Lawyer),
		HomeConfig: &HomeConfigService{repo: repo.HomeConfig, projectRepo: repo.Project, caseRepo: repo.Case, testimonialRepo: repo.Testimonial},
		Media: &MediaService{
			repo:            repo.Media,
			projectRepo:     repo.Project,
			caseRepo:        repo.Case,
			pageRepo:        repo.Page,
			lawyerRepo:      repo.Lawyer,
			testimonialRepo: repo.Testimonial,
			homeConfigRepo:  repo.HomeConfig,
		},
		Nav:           &NavService{repo: repo.Nav, projectRepo: repo.Project, pageRepo: repo.Page},
		Search:        &SearchService{faqRepo: repo.FAQ, pageRepo: repo.Page},
		Requirement:   &RequirementService{repo: repo.Requirement},
		CostItem:      &CostItemService{repo: repo.CostItem},
		TimelinePhase: &TimelinePhaseService{repo: repo.TimelinePhase},
		CompareConfig: &CompareConfigService{repo: repo.CompareConfig},
		Advantage:     &ProjectAdvantageService{repo: repo.ProjectAdvantage},
		Testimonial:   &TestimonialService{repo: repo.Testimonial},
		PublicVersion: NewPublicVersionService(),
	}

	// Wire home_config cleanup into entity services
	svc.Case.homeConfigSvc = svc.HomeConfig
	svc.Project.homeConfigSvc = svc.HomeConfig
	svc.Testimonial.homeConfigSvc = svc.HomeConfig

	// Wire projectRepo into LeadService for preloading project names
	svc.Lead.projectRepo = repo.Project

	// Post-wire cascade delete dependencies into ProjectService
	svc.Project.requirementRepo = repo.Requirement
	svc.Project.costItemRepo = repo.CostItem
	svc.Project.timelinePhaseRepo = repo.TimelinePhase
	svc.Project.milestoneRepo = repo.Milestone
	svc.Project.advantageRepo = repo.ProjectAdvantage
	svc.Project.compareConfigRepo = repo.CompareConfig
	svc.Project.caseRepo = repo.Case
	svc.Project.testimonialRepo = repo.Testimonial
	svc.Project.faqRepo = repo.FAQ
	svc.Project.versionRepo = repo.PublicVersion
	svc.Page.versionRepo = repo.PublicVersion
	svc.Case.versionRepo = repo.PublicVersion
	svc.FAQ.versionRepo = repo.PublicVersion
	svc.HomeConfig.versionRepo = repo.PublicVersion
	svc.Nav.versionRepo = repo.PublicVersion
	svc.Lawyer.versionRepo = repo.PublicVersion
	svc.Testimonial.versionRepo = repo.PublicVersion

	registerPublicVersionRegistrars(svc)

	return svc
}

func registerPublicVersionRegistrars(svc *Service) {
	v := reflect.ValueOf(svc).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).PkgPath != "" || t.Field(i).Name == "PublicVersion" {
			continue
		}
		if registrar, ok := v.Field(i).Interface().(PublicVersionRegistrar); ok {
			registrar.RegisterPublicVersions(svc.PublicVersion.registry)
		}
	}
}

// Stats returns dashboard statistics with month-over-month trends.
func (s *Service) Stats() (*dto.DashboardStats, error) {
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart

	projects, err := s.repo.Project.Count()
	if err != nil {
		logging.Logger.Warn("stats: failed to count projects", "error", err)
	}
	pages, err := s.repo.Page.Count()
	if err != nil {
		logging.Logger.Warn("stats: failed to count pages", "error", err)
	}
	leads, err := s.repo.Lead.Count()
	if err != nil {
		logging.Logger.Warn("stats: failed to count leads", "error", err)
	}
	cases, err := s.repo.Case.Count()
	if err != nil {
		logging.Logger.Warn("stats: failed to count cases", "error", err)
	}
	unread, err := s.repo.Lead.CountByStatus("new")
	if err != nil {
		logging.Logger.Warn("stats: failed to count unread leads", "error", err)
	}

	projectsThis, err := s.repo.Project.CountByRange(thisMonthStart, now)
	if err != nil {
		logging.Logger.Warn("stats: failed to count projects this month", "error", err)
	}
	projectsLast, err := s.repo.Project.CountByRange(lastMonthStart, lastMonthEnd)
	if err != nil {
		logging.Logger.Warn("stats: failed to count projects last month", "error", err)
	}
	pagesThis, err := s.repo.Page.CountByRange(thisMonthStart, now)
	if err != nil {
		logging.Logger.Warn("stats: failed to count pages this month", "error", err)
	}
	pagesLast, err := s.repo.Page.CountByRange(lastMonthStart, lastMonthEnd)
	if err != nil {
		logging.Logger.Warn("stats: failed to count pages last month", "error", err)
	}
	leadsThis, err := s.repo.Lead.CountByRange(thisMonthStart, now)
	if err != nil {
		logging.Logger.Warn("stats: failed to count leads this month", "error", err)
	}
	leadsLast, err := s.repo.Lead.CountByRange(lastMonthStart, lastMonthEnd)
	if err != nil {
		logging.Logger.Warn("stats: failed to count leads last month", "error", err)
	}
	casesThis, err := s.repo.Case.CountByRange(thisMonthStart, now)
	if err != nil {
		logging.Logger.Warn("stats: failed to count cases this month", "error", err)
	}
	casesLast, err := s.repo.Case.CountByRange(lastMonthStart, lastMonthEnd)
	if err != nil {
		logging.Logger.Warn("stats: failed to count cases last month", "error", err)
	}

	calc := func(this, last int64) dto.Trend {
		if last == 0 {
			if this == 0 {
				return dto.Trend{Direction: "neutral", Percent: 0, Label: "持平"}
			}
			return dto.Trend{Direction: "up", Percent: 100, Label: "本月新增"}
		}
		diff := float64(this-last) / float64(last) * 100
		if diff >= 0 {
			pct := int(diff + 0.5)
			if pct == 0 {
				return dto.Trend{Direction: "neutral", Percent: 0, Label: "持平"}
			}
			return dto.Trend{Direction: "up", Percent: pct, Label: fmtLabel("up", pct)}
		}
		pct := int(-diff + 0.5)
		return dto.Trend{Direction: "down", Percent: pct, Label: fmtLabel("down", pct)}
	}

	pT := calc(projectsThis, projectsLast)
	pgT := calc(pagesThis, pagesLast)
	lT := calc(leadsThis, leadsLast)
	cT := calc(casesThis, casesLast)

	return &dto.DashboardStats{
		TotalProjects: projects,
		TotalPages:    pages,
		TotalLeads:    leads,
		TotalCases:    cases,
		UnreadLeads:   unread,
		Trends: []dto.Trend{
			{Key: "projects", Direction: pT.Direction, Percent: pT.Percent, Label: pT.Label},
			{Key: "pages", Direction: pgT.Direction, Percent: pgT.Percent, Label: pgT.Label},
			{Key: "leads", Direction: lT.Direction, Percent: lT.Percent, Label: lT.Label},
			{Key: "cases", Direction: cT.Direction, Percent: cT.Percent, Label: cT.Label},
		},
	}, nil
}

func fmtLabel(dir string, pct int) string {
	if dir == "up" {
		return "较上月 +" + strconv.Itoa(pct) + "%"
	}
	return "较上月 -" + strconv.Itoa(pct) + "%"
}
