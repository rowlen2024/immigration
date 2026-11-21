package main

import (
	"testing"

	"mygo-immigration/backend/internal/model"
)

func TestPlanNavigationRegionsPreservesFooterAndHeader(t *testing.T) {
	root, wrapper, secondWrapper := uint64(1), uint64(2), uint64(8)
	project, second, latent, distinct := uint64(21), uint64(22), uint64(23), uint64(24)
	support := "/contact"
	rows := []model.Navigation{
		{ID: root, Status: true, DisplayPosition: "both"},
		{ID: wrapper, ParentID: &root, Label: "巴拿马", LinkType: "project", ProjectID: &project, Status: true, DisplayPosition: "both"},
		{ID: 3, ParentID: &wrapper, Label: "巴拿马项目", LinkType: "project", ProjectID: &project, Status: true, DisplayPosition: "header"},
		{ID: 4, ParentID: &wrapper, Label: "第二个项目", LinkType: "project", ProjectID: &second, Status: true, DisplayPosition: "both"},
		{ID: 5, ParentID: &wrapper, Label: "旧页脚没有渲染的项目", LinkType: "project", ProjectID: &latent, Status: true, DisplayPosition: "footer"},
		{ID: 6, ParentID: &wrapper, Label: "联系我们", LinkType: "custom", Link: &support, Status: true, DisplayPosition: "header"},
		{ID: secondWrapper, ParentID: &root, Label: "美国", LinkType: "project", ProjectID: &distinct, Status: true, DisplayPosition: "both"},
		{ID: 9, ParentID: &secondWrapper, Label: "其他地区分组重复的项目", LinkType: "project", ProjectID: &second, Status: true, DisplayPosition: "header"},
	}
	plan, err := planNavigationRegions(rows, []uint64{secondWrapper, wrapper}, map[uint64]string{distinct: "美国独立项目"})
	if err != nil {
		t.Fatal(err)
	}
	if plan.HeaderLinks != 4 || plan.FooterLinks != 2 {
		t.Fatalf("unexpected header/footer link counts: %+v", plan)
	}
	changed := make(map[uint64]model.Navigation)
	for _, c := range plan.Changes {
		changed[c.After.ID] = c.After
	}
	if changed[wrapper].Status || changed[5].Status || changed[9].Status {
		t.Fatalf("legacy wrapper, latent footer item and duplicate should be inactive: %+v", changed)
	}
	if promoted := changed[3]; promoted.ParentID == nil || *promoted.ParentID != root || promoted.DisplayPosition != "both" || promoted.Label != "巴拿马项目" || promoted.SortOrder != 0 {
		t.Fatalf("promoted project should preserve its label, order and both positions: %+v", promoted)
	}
	if changed[4].DisplayPosition != "header" || changed[6].Link == nil || *changed[6].Link != support {
		t.Fatalf("header-only links should be promoted without being added to the footer: %+v", changed)
	}
	if direct := changed[secondWrapper]; !direct.Status || direct.Label != "美国独立项目" || direct.DisplayPosition != "both" {
		t.Fatalf("wrapper with a separate project should become its direct link: %+v", direct)
	}
}

func TestPlanNavigationRegionsRejectsRootOrLeafSelection(t *testing.T) {
	root := uint64(1)
	rows := []model.Navigation{
		{ID: root, Status: true, DisplayPosition: "both"},
		{ID: 2, ParentID: &root, Status: true, DisplayPosition: "both"},
	}
	for _, id := range []uint64{1, 2, 3} {
		if _, err := planNavigationRegions(rows, []uint64{id}, nil); err == nil {
			t.Fatalf("expected navigation %d to be rejected", id)
		}
	}
}
