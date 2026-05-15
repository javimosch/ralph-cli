package gitops

import (
	"fmt"
	"os/exec"
	"strings"
)

// EnsureBranch creates the branch if it doesn't exist and checks it out.
func EnsureBranch(workDir, branchName string) error {
	// Check if branch already exists using git branch --list (refs only, not tags)
	cmd := exec.Command("git", "branch", "--list", branchName)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git branch list: %w", err)
	}

	if strings.TrimSpace(string(out)) == "" {
		// Branch doesn't exist, create it from current state
		createCmd := exec.Command("git", "checkout", "-b", branchName)
		createCmd.Dir = workDir
		if err := createCmd.Run(); err != nil {
			return fmt.Errorf("create branch %s: %w", branchName, err)
		}
		return nil
	}

	// Branch exists, check it out
	checkoutCmd := exec.Command("git", "checkout", branchName)
	checkoutCmd.Dir = workDir
	return checkoutCmd.Run()
}

// CommitAll stages all changes and commits them.
func CommitAll(workDir, message string) error {
	// Stage all changes
	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = workDir
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	// Check if there's anything to commit
	diffCmd := exec.Command("git", "diff", "--cached", "--quiet")
	diffCmd.Dir = workDir
	if diffCmd.Run() == nil {
		// No changes to commit
		return nil
	}

	// Commit
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit: %s: %w", strings.TrimSpace(string(out)), err)
	}

	return nil
}

// Status returns a summary of the current git state.
func Status(workDir string) (string, error) {
	cmd := exec.Command("git", "status", "--short", "--branch")
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git status: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CurrentBranch returns the name of the current git branch.
func CurrentBranch(workDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// IsRepo checks if workDir is inside a git repository.
func IsRepo(workDir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = workDir
	return cmd.Run() == nil
}
