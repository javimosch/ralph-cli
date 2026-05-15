package prd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javimosch/ralph-cli/util"
)

// Story represents a single user story / task in the PRD.
type Story struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	AcceptanceCriteria []string `json:"acceptanceCriteria"`
	Priority           int      `json:"priority"`
	Passes             bool     `json:"passes"`
	Notes              string   `json:"notes"`
	DependsOn          []string `json:"dependsOn"`
}

// PRD represents a Product Requirements Document file.
type PRD struct {
	Name        string  `json:"name"`
	BranchName  string  `json:"branchName,omitempty"`
	Description string  `json:"description,omitempty"`
	UserStories []Story `json:"userStories"`
}

// Load reads a PRD JSON file from the given path.
func Load(path string) (*PRD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read prd file: %w", err)
	}

	var prd PRD
	if err := json.Unmarshal(data, &prd); err != nil {
		return nil, fmt.Errorf("parse prd json: %w", err)
	}

	if err := prd.Validate(); err != nil {
		return nil, fmt.Errorf("validate prd: %w", err)
	}

	return &prd, nil
}

// Save writes the PRD back to its JSON file.
func (p *PRD) Save(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal prd: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	return os.WriteFile(path, append(data, '\n'), 0644)
}

// Validate checks that the PRD has all required fields and no dependency issues.
func (p *PRD) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("prd name is required")
	}
	if len(p.UserStories) == 0 {
		return fmt.Errorf("prd must have at least one user story")
	}
	seen := make(map[string]bool)
	for i, s := range p.UserStories {
		if s.ID == "" {
			return fmt.Errorf("story %d: id is required", i)
		}
		if s.Title == "" {
			return fmt.Errorf("story %s: title is required", s.ID)
		}
		if seen[s.ID] {
			return fmt.Errorf("duplicate story id: %s", s.ID)
		}
		seen[s.ID] = true
	}

	// Validate dependsOn references
	for _, s := range p.UserStories {
		for _, dep := range s.DependsOn {
			if !seen[dep] {
				return fmt.Errorf("story %s depends on unknown story: %s", s.ID, dep)
			}
		}
	}

	// Detect circular dependencies
	if cycle := p.findCycle(); cycle != nil {
		return fmt.Errorf("circular dependency detected: %v", cycle)
	}

	return nil
}

// findCycle detects circular dependencies between stories using DFS.
func (p *PRD) findCycle() []string {
	// Build adjacency map
	adj := make(map[string][]string)
	for _, s := range p.UserStories {
		adj[s.ID] = s.DependsOn
	}

	const (
		white = 0 // unvisited
		gray  = 1 // in current path
		black = 2 // done
	)

	color := make(map[string]int)
	for _, s := range p.UserStories {
		color[s.ID] = white
	}

	var path []string
	var dfs func(id string) bool
	dfs = func(id string) bool {
		color[id] = gray
		path = append(path, id)

		for _, dep := range adj[id] {
			if color[dep] == gray {
				// Found cycle - include dep in path
				path = append(path, dep)
				return true
			}
			if color[dep] == white {
				if dfs(dep) {
					return true
				}
			}
		}

		color[id] = black
		path = path[:len(path)-1]
		return false
	}

	for _, s := range p.UserStories {
		if color[s.ID] == white {
			path = path[:0]
			if dfs(s.ID) {
				return path
			}
		}
	}

	return nil
}

// WorkDir returns the working directory derived from the branch name or prd name.
func (p *PRD) WorkDir() string {
	if p.BranchName != "" {
		return p.BranchName
	}
	return fmt.Sprintf("ralph/%s", util.Slugify(p.Name))
}
