package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/microcosm-cc/bluemonday"
)

// LeadService 负责客户咨询的保存、查询和异步通知。
type LeadService struct {
	// repo 保存和查询咨询记录。
	repo repository.LeadRepository
	// projectRepo 提供咨询列表中展示的项目名称。
	projectRepo repository.ProjectRepository
	// notifier 处理成功保存后的咨询通知，构造时必须注入。
	notifier LeadNotifier
}

// NewLeadService 注入必需的咨询仓储和通知器。
func NewLeadService(repo repository.LeadRepository, notifier LeadNotifier) *LeadService {
	return &LeadService{
		repo:     repo,
		notifier: notifier,
	}
}

// leadSanitizer 清除咨询展示文本中的 HTML 标签。
var leadSanitizer = bluemonday.StrictPolicy()

// Create 净化并保存咨询，保存成功后异步通知且不等待邮件结果。
func (s *LeadService) Create(req *dto.LeadRequest) (*model.Lead, error) {
	if req == nil {
		return nil, errors.New("lead is nil")
	}
	if req.Name == "" {
		return nil, errors.New("lead name is required")
	}
	if req.Phone == "" {
		return nil, errors.New("lead phone is required")
	}

	lead := &model.Lead{
		Name:              leadSanitizer.Sanitize(req.Name),
		Phone:             req.Phone,
		Email:             req.Email,
		InterestedProject: req.InterestedProject,
		Message:           leadSanitizer.Sanitize(req.Message),
	}
	if err := s.repo.Create(lead); err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	// goroutine 参数在启动前按值复制，避免与返回记录的后续修改共享数据。
	go s.notifyLead(*lead)
	return lead, nil
}

// notifyLead 处理独立于 HTTP 请求的通知，失败只记录日志。
func (s *LeadService) notifyLead(lead model.Lead) {
	if err := s.notifier.Notify(lead); err != nil {
		logging.Logger.Error("failed to send lead notification", "lead_id", lead.ID, "error", err)
	}
}

// List returns leads with optional filters and pagination.
func (s *LeadService) List(req dto.LeadListRequest) ([]model.Lead, int64, error) {
	leads, total, err := s.repo.FindAll(repository.LeadFilter{
		Status:            req.Status,
		Name:              req.Name,
		Email:             req.Email,
		InterestedProject: req.InterestedProject,
		Page:              req.Page,
		PerPage:           req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", err)
	}

	// preload 项目名称
	if s.projectRepo != nil && len(leads) > 0 {
		slugSet := make(map[string]struct{})
		for _, l := range leads {
			if l.InterestedProject != "" {
				slugSet[l.InterestedProject] = struct{}{}
			}
		}
		if len(slugSet) > 0 {
			slugs := make([]string, 0, len(slugSet))
			for slug := range slugSet {
				slugs = append(slugs, slug)
			}
			projects, err := s.projectRepo.FindBySlugsLight(slugs)
			if err == nil {
				nameBySlug := make(map[string]string, len(projects))
				for _, p := range projects {
					nameBySlug[p.Slug] = p.Name
				}
				for i := range leads {
					if name, ok := nameBySlug[leads[i].InterestedProject]; ok {
						leads[i].ProjectName = name
					}
				}
			}
		}
	}

	return leads, total, nil
}

// Update updates the status and notes of a lead by ID.
func (s *LeadService) Update(id uint64, status, notes string) (*model.Lead, error) {
	if id == 0 {
		return nil, errors.New("lead id is required")
	}
	validStatuses := map[string]bool{
		"new": true, "contacted": true, "qualified": true, "closed": true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	if err := s.repo.UpdateStatus(id, status, notes); err != nil {
		return nil, fmt.Errorf("failed to update lead status: %w", err)
	}
	return &model.Lead{ID: id, Status: status, Notes: notes}, nil
}
