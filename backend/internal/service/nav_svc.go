package service

import (
	"fmt"
	"strings"
	"sync"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type NavService struct {
	repo        repository.NavigationRepository
	projectRepo repository.ProjectRepository
	pageRepo    repository.PageRepository
}

func (s *NavService) GetTree(position string) ([]model.Navigation, error) {
	items, err := s.repo.FindAllActiveByPosition(position)
	if err != nil {
		return nil, err
	}
	tree := buildTree(items, nil)
	s.fillLinks(tree)
	return tree, nil
}

func (s *NavService) GetAdminTree() ([]model.Navigation, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	tree := buildTree(items, nil)
	s.fillLinks(tree)
	return tree, nil
}

func (s *NavService) AdminList(page, pageSize int) ([]model.Navigation, int64, error) {
	allItems, err := s.repo.FindAll()
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(allItems))
	offset := (page - 1) * pageSize
	if offset >= int(total) {
		return []model.Navigation{}, total, nil
	}
	end := offset + pageSize
	if end > int(total) {
		end = int(total)
	}
	return allItems[offset:end], total, nil
}

func (s *NavService) Create(nav *model.Navigation) (*model.Navigation, error) {
	if err := s.validateNav(nav); err != nil {
		return nil, err
	}
	if nav.ParentID != nil {
		if _, err := s.repo.FindByID(*nav.ParentID); err != nil {
			return nil, fmt.Errorf("父导航项不存在")
		}
		d, err := s.getDepth(nav.ParentID)
		if err != nil {
			return nil, fmt.Errorf("父导航项不存在")
		}
		if d+1 > 3 {
			return nil, fmt.Errorf("导航最大层级为3级")
		}
	}
	if err := s.repo.Create(nav); err != nil {
		return nil, err
	}
	s.fillLink(nav)
	return nav, nil
}

func (s *NavService) Update(id uint64, nav *model.Navigation) (*model.Navigation, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("导航项不存在")
	}
	if err := s.validateNav(nav); err != nil {
		return nil, err
	}
	if nav.ParentID != nil {
		if *nav.ParentID == id {
			return nil, fmt.Errorf("不能将导航项设为自身的子项")
		}
		if _, err := s.repo.FindByID(*nav.ParentID); err != nil {
			return nil, fmt.Errorf("父导航项不存在")
		}
		if s.isDescendantOf(*nav.ParentID, id) {
			return nil, fmt.Errorf("不能将导航项设为其后代项的子项")
		}
		d, err := s.getDepth(nav.ParentID)
		if err != nil {
			return nil, fmt.Errorf("父导航项不存在")
		}
		if d+1 > 3 {
			return nil, fmt.Errorf("导航最大层级为3级")
		}
	}
	nav.ID = id
	nav.CreatedAt = existing.CreatedAt
	if err := s.repo.Update(nav); err != nil {
		return nil, err
	}
	s.fillLink(nav)
	return nav, nil
}

func (s *NavService) Delete(id uint64) error {
	hasChildren, err := s.repo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return fmt.Errorf("请先删除子导航项")
	}
	return s.repo.Delete(id)
}

func (s *NavService) validateNav(nav *model.Navigation) error {
	if strings.TrimSpace(nav.Label) == "" {
		return fmt.Errorf("导航名称不能为空")
	}
	if nav.LinkType == "" {
		nav.LinkType = "custom"
	}
	switch nav.LinkType {
	case "project":
		if nav.ProjectID == nil {
			return fmt.Errorf("项目链接必须选择项目")
		}
		nav.Link = nil // generated dynamically
	case "page":
		if nav.PageID == nil {
			return fmt.Errorf("页面链接必须选择页面")
		}
		nav.Link = nil // generated dynamically
	case "custom":
		if nav.Link == nil || strings.TrimSpace(*nav.Link) == "" {
			return fmt.Errorf("自定义链接不能为空")
		}
		if !isValidInternalURL(*nav.Link) {
			return fmt.Errorf("链接必须为内部路由（以/开头）")
		}
	default:
		return fmt.Errorf("不支持的链接类型: %s", nav.LinkType)
	}
	return nil
}

func (s *NavService) fillLink(nav *model.Navigation) {
	switch nav.LinkType {
	case "project":
		if nav.ProjectID != nil {
			slug := s.lookupProjectSlug(*nav.ProjectID)
			link := "/projects/" + slug
			nav.Link = &link
		}
	case "page":
		if nav.PageID != nil {
			slug := s.lookupPageSlug(*nav.PageID)
			link := "/" + slug
			nav.Link = &link
		}
	}
}

func (s *NavService) fillLinks(items []model.Navigation) {
	for i := range items {
		switch items[i].LinkType {
		case "project":
			if items[i].ProjectID != nil {
				slug := s.lookupProjectSlug(*items[i].ProjectID)
				link := "/projects/" + slug
				items[i].Link = &link
			}
		case "page":
			if items[i].PageID != nil {
				slug := s.lookupPageSlug(*items[i].PageID)
				link := "/" + slug
				items[i].Link = &link
			}
		}
		// custom: keep the Link as-is (already stored)
		if len(items[i].Children) > 0 {
			s.fillLinks(items[i].Children)
		}
	}
}

// slugCache avoids repeated DB lookups during tree traversal.
var (
	projectSlugMu   sync.RWMutex
	projectSlugCache = map[uint64]string{}
	pageSlugMu      sync.RWMutex
	pageSlugCache   = map[uint64]string{}
)

func (s *NavService) lookupProjectSlug(id uint64) string {
	projectSlugMu.RLock()
	if slug, ok := projectSlugCache[id]; ok {
		projectSlugMu.RUnlock()
		return slug
	}
	projectSlugMu.RUnlock()

	projects, _, _ := s.projectRepo.FindAll(1, 1000)

	projectSlugMu.Lock()
	for _, p := range projects {
		projectSlugCache[p.ID] = p.Slug
	}
	projectSlugMu.Unlock()

	for _, p := range projects {
		if p.ID == id {
			return p.Slug
		}
	}
	return ""
}

func (s *NavService) lookupPageSlug(id uint64) string {
	pageSlugMu.RLock()
	if slug, ok := pageSlugCache[id]; ok {
		pageSlugMu.RUnlock()
		return slug
	}
	pageSlugMu.RUnlock()

	pages, _ := s.pageRepo.FindAll()

	pageSlugMu.Lock()
	for _, p := range pages {
		pageSlugCache[p.ID] = p.Slug
	}
	pageSlugMu.Unlock()

	for _, p := range pages {
		if p.ID == id {
			return p.Slug
		}
	}
	return ""
}

func (s *NavService) CountByProjectID(projectID uint64) (int64, error) {
	return s.repo.CountByProjectID(projectID)
}

func (s *NavService) CountByPageID(pageID uint64) (int64, error) {
	return s.repo.CountByPageID(pageID)
}

func (s *NavService) getDepth(navID *uint64) (int, error) {
	if navID == nil {
		return 0, nil
	}
	depth := 0
	currentID := *navID
	for depth < 3 {
		nav, err := s.repo.FindByID(currentID)
		if err != nil {
			return 0, err
		}
		if nav.ParentID == nil {
			return depth + 1, nil
		}
		currentID = *nav.ParentID
		depth++
	}
	return depth, nil
}

func (s *NavService) isDescendantOf(parentID uint64, targetID uint64) bool {
	depth := 0
	currentID := parentID
	for depth < 3 {
		if currentID == targetID {
			return true
		}
		nav, err := s.repo.FindByID(currentID)
		if err != nil {
			return false
		}
		if nav.ParentID == nil {
			return false
		}
		currentID = *nav.ParentID
		depth++
	}
	return false
}

func isValidInternalURL(link string) bool {
	return strings.HasPrefix(link, "/") && !strings.Contains(link, "://")
}

func buildTree(items []model.Navigation, parentID *uint64) []model.Navigation {
	result := make([]model.Navigation, 0)
	for _, item := range items {
		matches := (parentID == nil && item.ParentID == nil) ||
			(parentID != nil && item.ParentID != nil && *item.ParentID == *parentID)
		if matches {
			item.Children = buildTree(items, &item.ID)
			result = append(result, item)
		}
	}
	return result
}
