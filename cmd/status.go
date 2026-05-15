package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javimosch/ralph-cli/gitops"
	"github.com/javimosch/ralph-cli/prd"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show PRD progress and current state",
	Long:  `Display the current progress of a PRD: completed vs total stories, blocked stories, and git branch info.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prdPath == "" {
			return fmt.Errorf("--prd flag is required")
		}

		prdFile, err := prd.Load(prdPath)
		if err != nil {
			return fmt.Errorf("load prd: %w", err)
		}

		completed, total := prdFile.Progress()
		next := prdFile.NextStory()
		blocked := prdFile.BlockedStories()

		workDir := "."
		branch, _ := gitops.CurrentBranch(workDir)
		isRepo := gitops.IsRepo(workDir)

		output := map[string]interface{}{
			"name":          prdFile.Name,
			"description":   prdFile.Description,
			"branch":        prdFile.WorkDir(),
			"currentBranch": branch,
			"isRepo":        isRepo,
			"completed":     completed,
			"total":         total,
			"isComplete":    prdFile.IsComplete(),
		}

		if next != nil {
			output["nextStory"] = map[string]interface{}{
				"id":    next.ID,
				"title": next.Title,
			}
		}

		if len(blocked) > 0 {
			blockedList := make([]map[string]interface{}, len(blocked))
			for i, s := range blocked {
				blockedList[i] = map[string]interface{}{
					"id":        s.ID,
					"title":     s.Title,
					"dependsOn": s.DependsOn,
				}
			}
			output["blockedStories"] = blockedList
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
