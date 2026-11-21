package main

import (
	"fmt"
	"reflect"
	"sort"

	"mygo-immigration/backend/internal/model"
)

const (
	navHeader = 1 << iota
	navFooter
)

type navRegionChange struct {
	Before model.Navigation `json:"before"`
	After  model.Navigation `json:"after"`
}

type navRegionPlan struct {
	WrapperIDs  []uint64          `json:"wrapper_ids"`
	HeaderLinks int               `json:"header_targets"`
	FooterLinks int               `json:"footer_targets"`
	Changes     []navRegionChange `json:"changes"`
}

func planNavigationRegions(rows []model.Navigation, wrapperIDs []uint64, names map[uint64]string) (navRegionPlan, error) {
	plan := navRegionPlan{WrapperIDs: append([]uint64(nil), wrapperIDs...), Changes: []navRegionChange{}}
	sort.Slice(plan.WrapperIDs, func(i, j int) bool { return plan.WrapperIDs[i] < plan.WrapperIDs[j] })
	if len(plan.WrapperIDs) == 0 {
		return plan, fmt.Errorf("请提供已审核的地区导航节点 ID")
	}
	before, children := indexNavigationRows(rows)
	after := make(map[uint64]model.Navigation, len(before))
	for id, row := range before {
		after[id] = row
	}
	wrappers := make(map[uint64]bool, len(wrapperIDs))
	roots := make(map[uint64]bool)
	for _, id := range plan.WrapperIDs {
		wrapper, ok := before[id]
		if !ok || !wrapper.Status || wrapper.ParentID == nil || len(children[id]) == 0 || wrappers[id] {
			return plan, fmt.Errorf("导航 %d 不是有效的待折叠地区节点", id)
		}
		root, ok := before[*wrapper.ParentID]
		if !ok || !root.Status || root.ParentID != nil {
			return plan, fmt.Errorf("导航 %d 必须是启用的顶级分组的直接子项", id)
		}
		for _, child := range children[id] {
			if len(children[child]) != 0 {
				return plan, fmt.Errorf("导航 %d 的子项 %d 仍有下级", id, child)
			}
		}
		wrappers[id] = true
		roots[root.ID] = true
	}

	for rootID := range roots {
		order := make([]uint64, 0, len(children[rootID]))
		moved := make(map[uint64]bool)
		for _, id := range children[rootID] {
			if !wrappers[id] {
				order = append(order, id)
				continue
			}
			wrapper := before[id]
			positions := navPositions(wrapper) & navPositions(before[rootID])
			var replacement uint64
			for _, childID := range children[id] {
				child := before[childID]
				child.ParentID = wrapper.ParentID
				// The old footer only rendered direct children of the footer column.
				setNavPositions(&child, navPositions(child)&positions&navHeader)
				if replacement == 0 && child.Status && navTarget(child) != "" && navTarget(child) == navTarget(wrapper) {
					replacement = childID
					setNavPositions(&child, navPositions(child)|positions)
				}
				after[childID] = child
				moved[childID] = true
			}

			switch {
			case replacement != 0:
				wrapper.Status = false
				order = append(order, replacement)
			case navTarget(wrapper) != "":
				if wrapper.LinkType == "project" && wrapper.ProjectID != nil && names[*wrapper.ProjectID] != "" {
					wrapper.Label = names[*wrapper.ProjectID]
				}
				moved[id] = true
				order = append(order, id)
			default:
				wrapper.Status = false
			}
			after[id] = wrapper
			for _, childID := range children[id] {
				if childID != replacement {
					order = append(order, childID)
				}
			}
		}
		dedupeNavigationProjects(order, children, moved, after)
		for rank, id := range order {
			row := after[id]
			row.SortOrder = rank
			after[id] = row
		}
	}

	afterRows := make([]model.Navigation, 0, len(rows))
	for _, old := range rows {
		row := after[old.ID]
		afterRows = append(afterRows, row)
		if !reflect.DeepEqual(old, row) {
			plan.Changes = append(plan.Changes, navRegionChange{Before: old, After: row})
		}
	}
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Before.ID < plan.Changes[j].Before.ID })
	for _, position := range []int{navHeader, navFooter} {
		previous := navigationTargets(rows, position)
		if !reflect.DeepEqual(previous, navigationTargets(afterRows, position)) {
			return plan, fmt.Errorf("折叠会改变 %s 的导航目标，未生成迁移计划", navPositionName(position))
		}
		if position == navHeader {
			plan.HeaderLinks = len(previous)
		} else {
			plan.FooterLinks = len(previous)
		}
	}
	return plan, nil
}

func indexNavigationRows(rows []model.Navigation) (map[uint64]model.Navigation, map[uint64][]uint64) {
	byID := make(map[uint64]model.Navigation, len(rows))
	children := make(map[uint64][]uint64)
	for _, row := range rows {
		byID[row.ID] = row
		if row.ParentID != nil {
			children[*row.ParentID] = append(children[*row.ParentID], row.ID)
		}
	}
	for _, ids := range children {
		sort.Slice(ids, func(i, j int) bool {
			left, right := byID[ids[i]], byID[ids[j]]
			if left.SortOrder == right.SortOrder {
				return left.ID < right.ID
			}
			return left.SortOrder < right.SortOrder
		})
	}
	return byID, children
}

func navPositions(row model.Navigation) int {
	if !row.Status {
		return 0
	}
	switch row.DisplayPosition {
	case "header":
		return navHeader
	case "footer":
		return navFooter
	case "both":
		return navHeader | navFooter
	default:
		return 0
	}
}

func navPositionName(positions int) string {
	switch positions {
	case navHeader:
		return "header"
	case navFooter:
		return "footer"
	default:
		return "both"
	}
}

func setNavPositions(row *model.Navigation, positions int) {
	row.Status = positions != 0
	if row.Status {
		row.DisplayPosition = navPositionName(positions)
	}
}

func navTarget(row model.Navigation) string {
	switch row.LinkType {
	case "project":
		if row.ProjectID != nil {
			return fmt.Sprintf("project:%d", *row.ProjectID)
		}
	case "page":
		if row.PageID != nil {
			return fmt.Sprintf("page:%d", *row.PageID)
		}
	default:
		if row.Link != nil {
			return *row.Link
		}
	}
	return ""
}

func dedupeNavigationProjects(order []uint64, children map[uint64][]uint64, moved map[uint64]bool, rows map[uint64]model.Navigation) {
	projects := make(map[uint64]uint64)
	for _, id := range order {
		row := rows[id]
		if !row.Status || row.LinkType != "project" || row.ProjectID == nil || (!moved[id] && len(children[id]) != 0) {
			continue
		}
		priorID := projects[*row.ProjectID]
		if priorID == 0 {
			projects[*row.ProjectID] = id
			continue
		}
		if moved[id] || moved[priorID] {
			prior := rows[priorID]
			setNavPositions(&prior, navPositions(prior)|navPositions(row))
			row.Status = false
			rows[priorID], rows[id] = prior, row
		}
	}
}

func navigationTargets(rows []model.Navigation, position int) map[string]bool {
	byID, children := indexNavigationRows(rows)
	targets := make(map[string]bool)
	var visit func(model.Navigation, uint64, int)
	visit = func(row model.Navigation, root uint64, depth int) {
		if navPositions(row)&position == 0 {
			return
		}
		if position == navHeader || depth == 1 {
			if target := navTarget(row); target != "" {
				targets[fmt.Sprintf("%d:%s", root, target)] = true
			}
		}
		if !((position == navFooter && depth == 1) || depth == 2) {
			for _, id := range children[row.ID] {
				visit(byID[id], root, depth+1)
			}
		}
	}
	for _, row := range rows {
		if row.ParentID == nil {
			visit(row, row.ID, 0)
		}
	}
	return targets
}
