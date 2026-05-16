package service

import (
	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/repository"
	"strconv"
	"time"
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
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo:          repo,
		Project:       &ProjectService{repo: repo.Project, navRepo: repo.Nav},
		Auth:          &AuthService{repo: repo.User, cfg: cfg},
		User:          &UserService{repo: repo.User},
		FAQ:           &FAQService{repo: repo.FAQ},
		Page:          &PageService{repo: repo.Page, navRepo: repo.Nav},
		Case:          &CaseService{repo: repo.Case},
		Lead:          &LeadService{repo: repo.Lead},
		Lawyer:        &LawyerService{repo: repo.Lawyer},
		HomeConfig:    &HomeConfigService{repo: repo.HomeConfig},
		Media:         &MediaService{repo: repo.Media},
		Nav:           &NavService{repo: repo.Nav, projectRepo: repo.Project, pageRepo: repo.Page},
		Search:        &SearchService{faqRepo: repo.FAQ, pageRepo: repo.Page},
		Requirement:   &RequirementService{repo: repo.Requirement},
		CostItem:      &CostItemService{repo: repo.CostItem},
		TimelinePhase: &TimelinePhaseService{repo: repo.TimelinePhase},
		CompareConfig: &CompareConfigService{repo: repo.CompareConfig},
		Advantage:     &ProjectAdvantageService{repo: repo.ProjectAdvantage},
	}
}

// DashboardStats holds admin dashboard statistics.
type DashboardStats struct {
	TotalProjects int64   `json:"totalProjects"`
	TotalPages    int64   `json:"totalPages"`
	TotalLeads    int64   `json:"totalLeads"`
	TotalCases    int64   `json:"totalCases"`
	UnreadLeads   int64   `json:"unreadLeads"`
	Trends        []Trend `json:"trends"`
}

// Trend represents a single metric's month-over-month trend.
type Trend struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`
	Percent   int    `json:"percent"`
	Label     string `json:"label"`
}

// Stats returns dashboard statistics with month-over-month trends.
func (s *Service) Stats() (*DashboardStats, error) {
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart

	projects, _ := s.repo.Project.Count()
	pages, _ := s.repo.Page.Count()
	leads, _ := s.repo.Lead.CountAll()
	cases, _ := s.repo.Case.Count()
	unread, _ := s.repo.Lead.CountByStatus("new")

	projectsThis, _ := s.repo.Project.CountByRange(thisMonthStart, now)
	projectsLast, _ := s.repo.Project.CountByRange(lastMonthStart, lastMonthEnd)
	pagesThis, _ := s.repo.Page.CountByRange(thisMonthStart, now)
	pagesLast, _ := s.repo.Page.CountByRange(lastMonthStart, lastMonthEnd)
	leadsThis, _ := s.repo.Lead.CountByRange(thisMonthStart, now)
	leadsLast, _ := s.repo.Lead.CountByRange(lastMonthStart, lastMonthEnd)
	casesThis, _ := s.repo.Case.CountByRange(thisMonthStart, now)
	casesLast, _ := s.repo.Case.CountByRange(lastMonthStart, lastMonthEnd)

	calc := func(this, last int64) Trend {
		if last == 0 {
			if this == 0 {
				return Trend{Direction: "neutral", Percent: 0, Label: "持平"}
			}
			return Trend{Direction: "up", Percent: 100, Label: "本月新增"}
		}
		diff := float64(this-last) / float64(last) * 100
		if diff >= 0 {
			pct := int(diff + 0.5)
			if pct == 0 {
				return Trend{Direction: "neutral", Percent: 0, Label: "持平"}
			}
			return Trend{Direction: "up", Percent: pct, Label: fmtLabel("up", pct)}
		}
		pct := int(-diff + 0.5)
		return Trend{Direction: "down", Percent: pct, Label: fmtLabel("down", pct)}
	}

	pT := calc(projectsThis, projectsLast)
	pgT := calc(pagesThis, pagesLast)
	lT := calc(leadsThis, leadsLast)
	cT := calc(casesThis, casesLast)

	return &DashboardStats{
		TotalProjects: projects,
		TotalPages:    pages,
		TotalLeads:    leads,
		TotalCases:    cases,
		UnreadLeads:   unread,
		Trends: []Trend{
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
