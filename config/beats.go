package config

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	beatsPath                           = flag.String("config", "./beats.yaml", "path to the beats config file")
	customParse                         = flag.Bool("custom-parse", false, "use custom config parse function")
	customParseFunc customParseFuncType = nil
	g               *BeatsConfig
)

type customParseFuncType func(p string) (*BeatsConfig, error)

func RegisterCustomParseFunc(f customParseFuncType) {
	if f != nil {
		customParseFunc = f
	}
}

type BeatsConfig struct {
	Beats          Beats             `yaml:"beats"`
	Tasks          []*TaskConfig     `yaml:"tasks"`
	CmdbMatchRules []*CmdbRuleConfig `yaml:"cmdb_match_rules,omitempty"`
}

type Beats struct {
	Logger    Logger `yaml:"logger"`
	TestModel bool   `yaml:"test_model"`
}

type Logger struct {
	Level   string `yaml:"level"`
	Output  string `yaml:"output"`
	LogFile string `yaml:"log_file"`
}

// InitConfig 进行配置的解析
func InitConfig() (*BeatsConfig, error) {
	if *customParse {
		return useCustomParseFunc(*beatsPath)
	}
	return defaultParseFunc(*beatsPath)
}

// useCustomParseFunc uses the custom parse function if it is registered.
func useCustomParseFunc(p string) (*BeatsConfig, error) {
	if customParseFunc != nil {
		return customParseFunc(p)
	}
	return nil, errors.New("custom parse function is not registered")
}

// defaultParseFunc parse config using the default method.
func defaultParseFunc(p string) (*BeatsConfig, error) {
	// 获取当前的工作目录
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 尝试进行配置文件路径拼接
	cPath := filepath.Join(dir, p)
	content, err := os.ReadFile(cPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content, &g); err != nil {
		return nil, err
	}
	return g, nil
}

func GetBeatsConfig() *BeatsConfig {
	return g
}

func GetLoggerConfig() *Logger {
	return &g.Beats.Logger
}

func GetTasks() []*TaskConfig {
	return g.Tasks
}
