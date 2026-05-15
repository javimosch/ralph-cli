package prd

import "sort"

// NextStory finds the highest-priority story that:
// - Has passes == false
// - All its dependencies (dependsOn) have passes == true
// Returns nil if all stories are complete or no story is available.
func (p *PRD) NextStory() *Story {
	available := make([]*Story, 0)

	for i := range p.UserStories {
		s := &p.UserStories[i]
		if s.Passes {
			continue
		}
		if depsBlocked(p.UserStories, s.DependsOn) {
			continue
		}
		available = append(available, s)
	}

	if len(available) == 0 {
		return nil
	}

	// Sort by priority (lower number = higher priority)
	sort.Slice(available, func(i, j int) bool {
		return available[i].Priority < available[j].Priority
	})

	return available[0]
}

// Progress returns counts of completed vs total stories.
func (p *PRD) Progress() (completed, total int) {
	total = len(p.UserStories)
	for _, s := range p.UserStories {
		if s.Passes {
			completed++
		}
	}
	return
}

// IsComplete returns true when all stories pass.
func (p *PRD) IsComplete() bool {
	for _, s := range p.UserStories {
		if !s.Passes {
			return false
		}
	}
	return true
}

// BlockedStories returns stories that are blocked by their dependencies.
func (p *PRD) BlockedStories() []Story {
	blocked := make([]Story, 0)
	for _, s := range p.UserStories {
		if !s.Passes && depsBlocked(p.UserStories, s.DependsOn) {
			blocked = append(blocked, s)
		}
	}
	return blocked
}

func depsBlocked(stories []Story, dependsOn []string) bool {
	if len(dependsOn) == 0 {
		return false
	}

	depsMap := make(map[string]bool, len(stories))
	for _, s := range stories {
		depsMap[s.ID] = s.Passes
	}

	for _, depID := range dependsOn {
		if !depsMap[depID] {
			return true
		}
	}
	return false
}
