package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javimosch/ralph-cli/template"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Scaffold a new PRD file",
	Long: `Create a new PRD JSON file with starter stories.
Provide a project name and optional description.
The output path can be set with --prd (default: ./prd.json).`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description := ""
		if len(args) > 1 {
			description = args[1]
		}

		outputPath := prdPath
		if outputPath == "" {
			outputPath = "./prd.json"
		}

		prd := template.Default(name, description)

		data, err := json.MarshalIndent(prd, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}

		if err := os.WriteFile(outputPath, append(data, '\n'), 0644); err != nil {
			return fmt.Errorf("write: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Created PRD at %s\n", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
