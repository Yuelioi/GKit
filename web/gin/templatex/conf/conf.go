package conf

import (
	"os"

	"github.com/goccy/go-yaml"
)

// Config 生成器配置
type Config struct {
	Model      ModelConfig      `yaml:"model"`
	API        APIConfig        `yaml:"api"`
	Database   DatabaseConfig   `yaml:"database"`
	Pagination PaginationConfig `yaml:"pagination"`
	Output     OutputConfig     `yaml:"output"`
	Features   FeaturesConfig   `yaml:"features"`
}

type ModelConfig struct {
	File       string   `yaml:"file"`        // 模型文件路径
	Package    string   `yaml:"package"`     // 模型包名
	Exclude    []string `yaml:"exclude"`     // 排除的模型
	Include    []string `yaml:"include"`     // 只包含的模型
	SoftDelete bool     `yaml:"soft_delete"` // 是否启用软删除
}

type APIConfig struct {
	Prefix     string            `yaml:"prefix"`      // API 路径前缀 如: /api/v1
	Group      string            `yaml:"group"`       // 路由组名
	Middleware []string          `yaml:"middleware"`  // 中间件
	Methods    map[string]bool   `yaml:"methods"`     // 启用的方法: create, get, list, update, delete
	CustomTags map[string]string `yaml:"custom_tags"` // 自定义 tag
}

type DatabaseConfig struct {
	Type       string `yaml:"type"`        // 数据库类型: mysql, postgres, sqlite
	TimeFields bool   `yaml:"time_fields"` // 是否自动处理时间字段
}

type PaginationConfig struct {
	Enabled     bool   `yaml:"enabled"`
	DefaultPage int    `yaml:"default_page"`
	DefaultSize int    `yaml:"default_size"`
	MaxSize     int    `yaml:"max_size"`
	Style       string `yaml:"style"` // offset 或 cursor
}

type OutputConfig struct {
	Dir       string            `yaml:"dir"`       // 输出目录
	Structure map[string]string `yaml:"structure"` // 目录结构自定义
	Overwrite bool              `yaml:"overwrite"` // 是否覆盖已存在文件
}

type FeaturesConfig struct {
	Validation   bool     `yaml:"validation"`    // 启用请求验证
	Cache        bool     `yaml:"cache"`         // 启用缓存
	Search       bool     `yaml:"search"`        // 启用搜索功能
	SearchFields []string `yaml:"search_fields"` // 可搜索字段
	Sort         bool     `yaml:"sort"`          // 启用排序
	Filter       bool     `yaml:"filter"`        // 启用过滤
	Export       bool     `yaml:"export"`        // 启用导出功能
}

type ModelInfo struct {
	Name         string
	Fields       []FieldInfo
	PkgName      string
	LowerName    string
	PluralName   string
	Config       *Config
	HasTimeField bool
}

type FieldInfo struct {
	Name       string
	Type       string
	JsonTag    string
	GormTag    string
	IsID       bool
	Required   bool
	Searchable bool
}

// 默认配置
func DefaultConfig() *Config {
	return &Config{
		Model: ModelConfig{
			SoftDelete: false,
		},
		API: APIConfig{
			Prefix: "/api",
			Methods: map[string]bool{
				"create": true,
				"get":    true,
				"list":   true,
				"update": true,
				"delete": true,
			},
		},
		Database: DatabaseConfig{
			Type:       "mysql",
			TimeFields: true,
		},
		Pagination: PaginationConfig{
			Enabled:     true,
			DefaultPage: 1,
			DefaultSize: 10,
			MaxSize:     100,
			Style:       "offset",
		},
		Output: OutputConfig{
			Dir: "generated",
			Structure: map[string]string{
				"handler": "handlers",
				"route":   "routes",
				"service": "services",
			},
			Overwrite: false,
		},
		Features: FeaturesConfig{
			Validation: true,
			Search:     false,
			Sort:       false,
			Filter:     false,
		},
	}
}

// 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	if configPath == "" {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
