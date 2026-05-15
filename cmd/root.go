package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	prdPath   string
	agentCmd  string
	dryRun    bool
	timeout   int
	noTimeout bool
)

var rootCmd = &cobra.Command{
	Use:   "ralph",
	Short: "Ralph — autonomous PRD-driven agent loop",
	Long: `Ralph is a minimal, file-based agent loop for autonomous coding.
It reads stories from a PRD JSON file and executes them one at a time,
using files and git as memory for fresh iterations and persistent state.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if noTimeout {
			timeout = 0
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&prdPath, "prd", "", "Path to PRD JSON file")
	rootCmd.PersistentFlags().StringVar(&agentCmd, "agent", "", "Agent CLI command to use (default: auto-detect)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print prompts without executing")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 300, "Max seconds per agent execution (default: 300)")
	rootCmd.PersistentFlags().BoolVar(&noTimeout, "no-timeout", false, "Disable execution timeout (overrides --timeout, runs until agent finishes)")
}
