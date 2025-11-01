package model

type Pipeline struct {
	ID     string   `json:"id" yaml:"id"`
	Name   string   `json:"name" yaml:"name"`
	Agent  *Agent   `json:"agent,omitempty" yaml:"agent,omitempty"`  // Pipeline-level agent
	Stages []Stage  `json:"stages" yaml:"stages"`
	Labels []string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type Agent struct {
	Label string `json:"label" yaml:"label"`
	// Add more configuration fields as needed
}

// StageType: "serial" (default) or "parallel"
type Stage struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type,omitempty" yaml:"type,omitempty"` // "", "serial", "parallel"
	Agent    *Agent   `json:"agent,omitempty" yaml:"agent,omitempty"` // Stage-level agent (overrides pipeline agent)
	Steps    []Step   `json:"steps,omitempty" yaml:"steps,omitempty"` // Used if serial stage
	Parallel []Stage  `json:"parallel,omitempty" yaml:"parallel,omitempty"` // Used if parallel stage
}

type Step struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"` // sh, py, go, scm, etc.
	Args     []string `json:"args,omitempty" yaml:"args,omitempty"`
	Script   string   `json:"script,omitempty" yaml:"script,omitempty"`
	Parallel []Step   `json:"parallel,omitempty" yaml:"parallel,omitempty"`
}