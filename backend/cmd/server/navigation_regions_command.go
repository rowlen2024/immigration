package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/database"
	"mygo-immigration/backend/internal/model"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func init() {
	var ids, applyFile, rollbackFile string
	cmd := &cobra.Command{
		Use:   "flatten-navigation-regions",
		Short: "预演、应用或回滚经审核的地区导航折叠",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runNavigationRegionCommand(cmd, ids, applyFile, rollbackFile)
		},
	}
	cmd.Flags().StringVar(&ids, "ids", "", "经审核的旧地区导航节点 ID，逗号分隔。只生成预演 JSON")
	cmd.Flags().StringVar(&applyFile, "apply", "", "应用已审核的预演 JSON")
	cmd.Flags().StringVar(&rollbackFile, "rollback", "", "按原预演 JSON 回滚已应用的变更")
	cmd.MarkFlagsOneRequired("ids", "apply", "rollback")
	cmd.MarkFlagsMutuallyExclusive("ids", "apply", "rollback")
	rootCmd.AddCommand(cmd)
}

func runNavigationRegionCommand(cmd *cobra.Command, ids, applyFile, rollbackFile string) error {
	var plan navRegionPlan
	if ids != "" {
		for _, value := range strings.Split(ids, ",") {
			id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
			if err != nil || id == 0 {
				return fmt.Errorf("无效的地区节点 ID %q", value)
			}
			plan.WrapperIDs = append(plan.WrapperIDs, id)
		}
	} else {
		planFile := applyFile
		if rollbackFile != "" {
			planFile = rollbackFile
		}
		data, err := os.ReadFile(planFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &plan); err != nil {
			return err
		}
		if len(plan.WrapperIDs) == 0 || len(plan.Changes) == 0 {
			return fmt.Errorf("预演文件没有可执行的导航变更")
		}
	}

	db, err := database.InitMySQL(config.Load())
	if err != nil {
		return err
	}
	db = db.Session(&gorm.Session{Context: cmd.Context(), Logger: logger.Default.LogMode(logger.Silent)})
	if ids != "" {
		preview, err := getNavigationRegionPlan(db, plan.WrapperIDs)
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		return encoder.Encode(preview)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		locked := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Session(&gorm.Session{})
		if applyFile != "" {
			current, err := getNavigationRegionPlan(locked, plan.WrapperIDs)
			if err != nil {
				return err
			}
			if !sameNavigationJSON(current, plan) {
				return fmt.Errorf("导航数据已变更，请重新预演")
			}
		}
		for _, change := range plan.Changes {
			from, to := change.Before, change.After
			if rollbackFile != "" {
				from, to = to, from
			}
			if from.ID == 0 || from.ID != to.ID {
				return fmt.Errorf("预演文件包含无效的导航变更")
			}
			if rollbackFile != "" {
				var current model.Navigation
				if err := locked.First(&current, from.ID).Error; err != nil {
					return err
				}
				current.UpdatedAt, from.UpdatedAt = time.Time{}, time.Time{}
				if !sameNavigationJSON(current, from) {
					return fmt.Errorf("导航 %d 上线后已被编辑，未回滚任何导航项", from.ID)
				}
			}
			if err := tx.Model(&model.Navigation{}).Where("id = ?", to.ID).Updates(map[string]any{
				"parent_id": to.ParentID, "label": to.Label, "sort_order": to.SortOrder,
				"status": to.Status, "display_position": to.DisplayPosition,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	cmd.Printf("已处理 %d 条导航记录\n", len(plan.Changes))
	return nil
}

func getNavigationRegionPlan(db *gorm.DB, ids []uint64) (navRegionPlan, error) {
	var rows []model.Navigation
	if err := db.Find(&rows).Error; err != nil {
		return navRegionPlan{}, err
	}
	selected := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		selected[id] = true
	}
	var projectIDs []uint64
	for _, row := range rows {
		if selected[row.ID] && row.LinkType == "project" && row.ProjectID != nil {
			projectIDs = append(projectIDs, *row.ProjectID)
		}
	}
	names := make(map[uint64]string, len(projectIDs))
	if len(projectIDs) > 0 {
		var projects []model.Project
		if err := db.Select("id", "name").Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
			return navRegionPlan{}, err
		}
		for _, project := range projects {
			names[project.ID] = project.Name
		}
	}
	return planNavigationRegions(rows, ids, names)
}

func sameNavigationJSON(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}
