package template

import (
	"github.com/javimosch/ralph-cli/prd"
	"github.com/javimosch/ralph-cli/util"
)

// Default returns a starter PRD with example stories.
func Default(name, description string) *prd.PRD {
	return &prd.PRD{
		Name:        name,
		BranchName:  util.Slugify(name),
		Description: description,
		UserStories: []prd.Story{
			{
				ID:          "US-001",
				Title:       "Project scaffolding and setup",
				Description: "As a developer, I want the project scaffolded so I can start implementing.",
				AcceptanceCriteria: []string{
					"Project structure is created",
					"Build system is configured",
					"Initial commit is made",
				},
				Priority: 1,
				Passes:   false,
				Notes:    "",
				DependsOn: []string{},
			},
		},
	}
}
