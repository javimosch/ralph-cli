package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javimosch/ralph-cli/loop"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the PRD agent loop",
	Long: `Run the autonomous agent loop against a PRD file.
Selects the next available story, spawns the agent, and commits results.
Repeat until all stories pass.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prdPath == "" {
			return fmt.Errorf("--prd flag is required")
		}

		results, err := loop.Run(loop.Opts{
			PRDPath:   prdPath,
			AgentCmd:  agentCmd,
			DryRun:    dryRun,
			Timeout:   timeout,
		})
		if err != nil {
			return fmt.Errorf("ralph loop: %w", err)
		}

		// Output results as JSON
		output := map[string]interface{}{
			"results": results,
			"status":  "completed",
		}

		// Check if any failed
		for _, r := range results {
			if r.Status == "failed" {
				output["status"] = "failed"
				break
			}
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
