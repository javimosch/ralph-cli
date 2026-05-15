package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javimosch/ralph-cli/loop"
	"github.com/javimosch/ralph-cli/prd"
	"github.com/spf13/cobra"
)

var storyID string

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Show or select the next story to work on",
}

var storyNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next available story",
	Long:  `Print the highest-priority story that is ready to work on (not blocked by dependencies).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prdPath == "" {
			return fmt.Errorf("--prd flag is required")
		}

		prdFile, err := prd.Load(prdPath)
		if err != nil {
			return fmt.Errorf("load prd: %w", err)
		}

		story := prdFile.NextStory()
		if story == nil {
			if prdFile.IsComplete() {
				fmt.Fprintln(os.Stderr, "All stories are complete!")
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, "No stories available - all remaining stories are blocked by dependencies.")
			os.Exit(1)
		}

		output := map[string]interface{}{
			"story":  story,
			"prompt": loop.BuildPrompt(prdFile, story),
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
	},
}

var storyPromptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Print the agent prompt for the next story",
	Long:  `Generate and print the agent prompt for the next available story without executing anything.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prdPath == "" {
			return fmt.Errorf("--prd flag is required")
		}

		prdFile, err := prd.Load(prdPath)
		if err != nil {
			return fmt.Errorf("load prd: %w", err)
		}

		// If storyID flag provided, show prompt for specific story
		var target *prd.Story
		if storyID != "" {
			for i := range prdFile.UserStories {
				if prdFile.UserStories[i].ID == storyID {
					target = &prdFile.UserStories[i]
					break
				}
			}
			if target == nil {
				return fmt.Errorf("story %s not found", storyID)
			}
		} else {
			target = prdFile.NextStory()
			if target == nil {
				if prdFile.IsComplete() {
					return fmt.Errorf("all stories are complete")
				}
				return fmt.Errorf("no available stories (blocked by dependencies)")
			}
		}

		fmt.Print(loop.BuildPrompt(prdFile, target))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storyCmd)
	storyCmd.AddCommand(storyNextCmd)
	storyCmd.AddCommand(storyPromptCmd)
	storyPromptCmd.Flags().StringVar(&storyID, "story", "", "Specific story ID (optional, defaults to next)")
}
