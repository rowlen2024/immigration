package service

import (
	"errors"
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"gorm.io/gorm"
)

func (s *ProjectService) PublishedProjectRegions(ids []uint64) ([]dto.PublishedProjectRegion, error) {
	projects, err := s.regionRepo.FindPublishedProjectRegions(ids)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PublishedProjectRegion, 0, len(projects))
	for _, project := range projects {
		item := dto.PublishedProjectRegion{ProjectID: project.ID}
		if project.Region != nil {
			region := projectRegionInfo(*project.Region)
			item.Region = &region
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *ProjectService) RegionOptions() ([]dto.ProjectRegionOption, error) {
	regions, err := s.regionRepo.ListRegions()
	if err != nil {
		return nil, err
	}
	result := make([]dto.ProjectRegionOption, 0, len(regions))
	for _, region := range regions {
		result = append(result, dto.ProjectRegionOption{ID: region.ID, ProjectRegion: projectRegionInfo(region)})
	}
	return result, nil
}

func projectRegionInfo(region model.ProjectRegion) dto.ProjectRegion {
	return dto.ProjectRegion{Code: region.Code, Name: region.Name, FlagEmoji: region.FlagEmoji}
}

func (s *ProjectService) publicProjectRegionsVersion(string) (repository.PublicVersion, error) {
	const query = `
SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM (
  SELECT updated_at FROM projects WHERE status = 1 AND deleted_at IS NULL
  UNION ALL SELECT updated_at FROM project_regions
) AS versions`
	return s.versionRepo.VersionFromQuery(query)
}

func (s *ProjectService) applyProjectRegion(project *model.Project, regionID *uint64, country string) error {
	if regionID != nil {
		region, err := s.regionRepo.FindRegionByID(*regionID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("所选国家或地区不存在")
		}
		if err != nil {
			return fmt.Errorf("failed to get project region: %w", err)
		}
		if project.ID == 0 || project.RegionID == nil || *project.RegionID != *regionID || project.Country == "" {
			project.Country = region.Name
		}
		project.RegionID = regionID
		return nil
	}

	if project.ID != 0 && project.Country == country {
		return nil
	}
	project.Country = country
	project.RegionID = nil
	name := strings.TrimSpace(country)
	if name == "" || s.regionRepo == nil {
		return nil
	}
	if name == "中国（香港特别行政区）" {
		name = "中国香港"
	}
	// 存量香港人才项目的「中国」需由运营选择，不从歧义文案推断。
	if project.Slug == "xiang-gang-gao-duan-ren-cai" && name == "中国" {
		return nil
	}
	region, err := s.regionRepo.FindRegionByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to match project region: %w", err)
	}
	project.RegionID = &region.ID
	return nil
}
