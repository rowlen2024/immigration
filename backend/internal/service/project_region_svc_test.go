package service

import (
	"testing"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"gorm.io/gorm"
)

type projectRegionLookupStub struct {
	repository.ProjectRegionRepository
}

func (projectRegionLookupStub) FindRegionByID(id uint64) (*model.ProjectRegion, error) {
	if id == 3 {
		return &model.ProjectRegion{ID: 3, Code: "HK", Name: "中国香港"}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (lookup projectRegionLookupStub) FindRegionByName(name string) (*model.ProjectRegion, error) {
	if name == "中国香港" {
		return lookup.FindRegionByID(3)
	}
	return nil, gorm.ErrRecordNotFound
}

func TestApplyProjectRegionKeepsCountryAndAssociationInSync(t *testing.T) {
	hongKong := uint64(3)
	usa := uint64(1)
	missing := uint64(42)
	const oldHongKongName = "中国（香港特别行政区）"
	svc := &ProjectService{regionRepo: projectRegionLookupStub{}}
	tests := []struct {
		name        string
		project     model.Project
		regionID    *uint64
		country     string
		wantCountry string
		wantRegion  *uint64
		wantError   bool
	}{
		{
			name: "same selection preserves legacy filter value", project: model.Project{ID: 5, Country: oldHongKongName, RegionID: &hongKong},
			regionID: &hongKong, country: "中国香港", wantCountry: oldHongKongName, wantRegion: &hongKong,
		},
		{
			name: "new selection uses standard name", project: model.Project{ID: 5, Country: "美国", RegionID: &usa},
			regionID: &hongKong, country: "美国", wantCountry: "中国香港", wantRegion: &hongKong,
		},
		{
			name: "old client resolves known alias", project: model.Project{ID: 5, Country: "美国", RegionID: &usa},
			country: oldHongKongName, wantCountry: oldHongKongName, wantRegion: &hongKong,
		},
		{
			name: "old client does not keep stale association", project: model.Project{ID: 5, Country: "香港", RegionID: &hongKong},
			country: "待确认地区", wantCountry: "待确认地区",
		},
		{
			name: "unknown selection rejected", project: model.Project{ID: 5, Country: "美国", RegionID: &usa},
			regionID: &missing, country: "待确认地区", wantCountry: "美国", wantRegion: &usa, wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := svc.applyProjectRegion(&test.project, test.regionID, test.country)
			if (err != nil) != test.wantError {
				t.Fatalf("applyProjectRegion error = %v, want error %t", err, test.wantError)
			}
			if test.project.Country != test.wantCountry {
				t.Errorf("country = %q, want %q", test.project.Country, test.wantCountry)
			}
			if (test.project.RegionID == nil) != (test.wantRegion == nil) || test.wantRegion != nil && *test.project.RegionID != *test.wantRegion {
				t.Error("region association does not match the expected region")
			}
		})
	}
}
