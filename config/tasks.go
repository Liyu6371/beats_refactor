package config

type TaskConfig struct {
	Name           string          `yaml:"name"`
	DataId         int64           `yaml:"data_id"`
	Period         string          `yaml:"period,omitempty"`
	Source         []*SourceConfig `yaml:"source,omitempty"`
	Pipeline       Pipeline        `yaml:"pipeline"`
	Sender         []*SenderConfig `yaml:"sender"`
	CmdbMatchRules []CmdbMatchRule `yaml:"cmdb_match_rules,omitempty"`
}

type CmdbMatchRule map[string]interface{}

type Pipeline struct {
	Processor []string `yaml:"processor"`
	Shaper    []string `yaml:"shaper"`
}

// IsCloudMonitorTask checks if the task is a cloud monitor task.
// A cloud monitor task is defined as one that which Source field is empty.
func (t *TaskConfig) IsCloudMonitorTask() bool {
	return len(t.Source) == 0 && t.Period == ""
}

type SourceConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type SenderConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"` // Configuration for the sender
}
