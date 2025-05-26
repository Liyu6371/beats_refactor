package config

type CmdbRuleConfig struct {
	Key         string       `yaml:"key"`
	Operator    string       `yaml:"operator"`
	Value       string       `yaml:"value"`
	ObjectModel string       `yaml:"object_model"`
	Rules       []*MatchRule `yaml:"instance_match_rules"`
}

type MatchRule struct {
	Key      string `yaml:"key"`
	Value    string `yaml:"value"`
	Operator string `yaml:"operator"`
}
