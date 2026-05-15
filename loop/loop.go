package loop

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/javimosch/ralph-cli/gitops"
	"github.com/javimosch/ralph-cli/prd"
)

type Opts struct {
	PRDPath   string
	AgentCmd  string
	AgentArgs []string
	DryRun    bool
	WorkDir   string
	Timeout   int // max seconds per agent execution (0 = no timeout)
}

type Result struct {
	StoryID string `json:"storyId"`
	Title   string `json:"title"`
	Status  string `json:"status"` // "passed", "skipped", "failed", "timedout"
	Error   string `json:"error,omitempty"`
}

// Run executes the ralph agent loop: one story at a time.
func Run(opts Opts) ([]Result, error) {
	prdFile, err := prd.Load(opts.PRDPath)
	if err != nil {
		return nil, fmt.Errorf("load prd: %w", err)
	}

	workDir := opts.WorkDir
	if workDir == "" {
		workDir = "."
	}

	results := make([]Result, 0)

	for !prdFile.IsComplete() {
		story := prdFile.NextStory()
		if story == nil {
			blocked := prdFile.BlockedStories()
			if len(blocked) > 0 {
				return results, fmt.Errorf("stories blocked by dependencies: %v", blockedIDs(blocked))
			}
			break
		}

		fmt.Fprintf(os.Stderr, "\n=== Ralph: Working on %s: %s ===\n", story.ID, story.Title)

		if opts.DryRun {
			fmt.Fprintf(os.Stderr, "[DRY RUN] Would execute story %s: %s\n", story.ID, story.Title)
			fmt.Fprintf(os.Stderr, "[DRY RUN] Prompt:\n%s\n", BuildPrompt(prdFile, story))
			results = append(results, Result{
				StoryID: story.ID,
				Title:   story.Title,
				Status:  "skipped",
			})
			story.Passes = true
			continue
		}

		// Create git branch for this PRD
		branchName := prdFile.WorkDir()
		if err := gitops.EnsureBranch(workDir, branchName); err != nil {
			return results, fmt.Errorf("git branch: %w", err)
		}

		// Generate and execute the agent prompt
		prompt := BuildPrompt(prdFile, story)
		res := runAgentWithTimeout(workDir, opts.AgentCmd, opts.AgentArgs, prompt, opts.Timeout)
		if res != nil {
			results = append(results, *res)
			return results, fmt.Errorf("agent execution for %s: %s", story.ID, res.Error)
		}

		// Mark story as passed
		story.Passes = true

		// Save updated prd.json
		if err := prdFile.Save(opts.PRDPath); err != nil {
			return results, fmt.Errorf("save prd after %s: %w", story.ID, err)
		}

		// Commit changes via git
		msg := fmt.Sprintf("ralph: %s - %s", story.ID, story.Title)
		if err := gitops.CommitAll(workDir, msg); err != nil {
			return results, fmt.Errorf("git commit for %s: %w", story.ID, err)
		}

		fmt.Fprintf(os.Stderr, "\n=== Ralph: %s passed! ===\n", story.ID)

		results = append(results, Result{
			StoryID: story.ID,
			Title:   story.Title,
			Status:  "passed",
		})
	}

	return results, nil
}

// BuildPrompt generates the prompt for an agent to work on a story.
func BuildPrompt(prdFile *prd.PRD, story *prd.Story) string {
	var b strings.Builder

	b.WriteString("You are working on a PRD-driven implementation.\n\n")
	fmt.Fprintf(&b, "PRD: %s\n", prdFile.Name)
	if prdFile.Description != "" {
		fmt.Fprintf(&b, "Description: %s\n", prdFile.Description)
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "=== Current Story: %s: %s ===\n\n", story.ID, story.Title)
	if story.Description != "" {
		fmt.Fprintf(&b, "Description:\n%s\n\n", story.Description)
	}

	b.WriteString("Acceptance Criteria:\n")
	for _, ac := range story.AcceptanceCriteria {
		fmt.Fprintf(&b, "- [ ] %s\n", ac)
	}
	b.WriteString("\n")

	if len(story.DependsOn) > 0 {
		b.WriteString("Prerequisites (already completed):\n")
		for _, depID := range story.DependsOn {
			for _, s := range prdFile.UserStories {
				if s.ID == depID && s.Passes {
					fmt.Fprintf(&b, "- %s: %s\n", s.ID, s.Title)
				}
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("Workflow:\n")
	b.WriteString("1. Read and understand the existing codebase\n")
	b.WriteString("2. Implement the changes needed to satisfy ALL acceptance criteria\n")
	b.WriteString("3. Verify your changes meet all criteria\n")
	b.WriteString("4. Report what you did\n")

	return b.String()
}

func runAgentWithTimeout(workDir, cmdName string, args []string, prompt string, timeoutSec int) *Result {
	if cmdName == "" {
		cmdName = detectDefaultAgent()
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if timeoutSec > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = workDir
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return &Result{
				Status: "timedout",
				Error:  fmt.Sprintf("timed out after %ds running %s", timeoutSec, cmdName),
			}
		}
		// Non-zero exit or other execution error
		return &Result{
			Status: "failed",
			Error:  fmt.Sprintf("%s exited with error: %s", cmdName, err.Error()),
		}
	}

	return nil // success
}

func detectDefaultAgent() string {
	candidates := []string{"opencode", "claude", "codex", "gemini-cli", "cursor"}
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}
	return "<agent-cli>"
}

func blockedIDs(stories []prd.Story) []string {
	ids := make([]string, len(stories))
	for i, s := range stories {
		ids[i] = s.ID
	}
	return ids
}
