package conf

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Model      ModelConfig      `yaml:"model"`
	API        APIConfig        `yaml:"api"`
	Database   DatabaseConfig   `yaml:"database"`
	Pagination PaginationConfig `yaml:"pagination"`
	Output     OutputConfig     `yaml:"output"`
	Features   FeaturesConfig   `yaml:"features"`
	Hooks      HooksConfig      `yaml:"hooks"`     // 新增
	Relations  RelationsConfig  `yaml:"relations"` // 新增
	Auth       AuthConfig       `yaml:"auth"`      // 新增
}

type ModelConfig struct {
	Files      []string          `yaml:"files"`
	Package    string            `yaml:"package"`
	Exclude    []string          `yaml:"exclude"`
	Include    []string          `yaml:"include"`
	SoftDelete bool              `yaml:"soft_delete"`
	Timestamps bool              `yaml:"timestamps"` // 新增：自动时间戳
	TableName  map[string]string `yaml:"table_name"` // 新增：自定义表名
}

type APIConfig struct {
	Prefix     string            `yaml:"prefix"`
	Group      string            `yaml:"group"`
	Middleware []string          `yaml:"middleware"`
	Methods    map[string]bool   `yaml:"methods"`
	CustomTags map[string]string `yaml:"custom_tags"`
	RateLimit  RateLimitConfig   `yaml:"rate_limit"` // 新增：限流
	CORS       CORSConfig        `yaml:"cors"`       // 新增：跨域
	Swagger    SwaggerConfig     `yaml:"swagger"`    // 新增：API文档
	Versioning bool              `yaml:"versioning"` // 新增：API版本控制
}

type RateLimitConfig struct {
	Enabled bool   `yaml:"enabled"`
	Rate    string `yaml:"rate"` // 如: "100/1m", "1000/1h"
	Burst   int    `yaml:"burst"`
}

type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type SwaggerConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Host        string `yaml:"host"`
	BasePath    string `yaml:"base_path"`
}

type DatabaseConfig struct {
	Type        string        `yaml:"type"`
	TimeFields  bool          `yaml:"time_fields"`
	Transaction bool          `yaml:"transaction"` // 新增：事务支持
	Preload     []string      `yaml:"preload"`     // 新增：预加载关联
	Indexes     []IndexConfig `yaml:"indexes"`     // 新增：索引配置
	Hooks       []string      `yaml:"hooks"`       // 新增：GORM钩子
}

type IndexConfig struct {
	Name   string   `yaml:"name"`
	Fields []string `yaml:"fields"`
	Unique bool     `yaml:"unique"`
}

type PaginationConfig struct {
	Enabled     bool   `yaml:"enabled"`
	DefaultPage int    `yaml:"default_page"`
	DefaultSize int    `yaml:"default_size"`
	MaxSize     int    `yaml:"max_size"`
	Style       string `yaml:"style"`
	CursorField string `yaml:"cursor_field"` // 新增：游标字段
}

type OutputConfig struct {
	Dir       string            `yaml:"dir"`
	Structure map[string]string `yaml:"structure"`
	Overwrite bool              `yaml:"overwrite"`
	Format    FormatConfig      `yaml:"format"`  // 新增：格式化
	License   string            `yaml:"license"` // 新增：许可证头
	Author    string            `yaml:"author"`  // 新增：作者信息
}

type FormatConfig struct {
	Enabled   bool `yaml:"enabled"`
	GoFmt     bool `yaml:"gofmt"`
	GoImports bool `yaml:"goimports"`
	Golangci  bool `yaml:"golangci"`
}

type FeaturesConfig struct {
	Validation   bool            `yaml:"validation"`
	Cache        CacheConfig     `yaml:"cache"` // 增强：详细缓存配置
	Search       bool            `yaml:"search"`
	SearchFields []string        `yaml:"search_fields"`
	Sort         bool            `yaml:"sort"`
	Filter       FilterConfig    `yaml:"filter"`     // 增强：过滤配置
	Export       ExportConfig    `yaml:"export"`     // 增强：导出配置
	Import       bool            `yaml:"import"`     // 新增：导入功能
	Batch        bool            `yaml:"batch"`      // 新增：批量操作
	Service      bool            `yaml:"service"`    // 新增：Service层
	Repository   bool            `yaml:"repository"` // 新增：Repository层
	DTO          bool            `yaml:"dto"`        // 新增：DTO层
	Logging      LoggingConfig   `yaml:"logging"`    // 新增：日志
	Metrics      bool            `yaml:"metrics"`    // 新增：监控指标
	Tracing      bool            `yaml:"tracing"`    // 新增：链路追踪
	Recovery     bool            `yaml:"recovery"`   // 新增：panic恢复
	Validator    ValidatorConfig `yaml:"validator"`  // 新增：验证器配置
	Test         bool            `yaml:"test"`       // 新增：测试生成
	Mock         bool            `yaml:"mock"`       // 新增：Mock生成
}

type CacheConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Type      string `yaml:"type"` // redis, memory, memcached
	TTL       string `yaml:"ttl"`  // 过期时间 如: "5m", "1h"
	KeyPrefix string `yaml:"key_prefix"`
	Strategy  string `yaml:"strategy"` // write-through, write-back, cache-aside
}

type FilterConfig struct {
	Enabled       bool     `yaml:"enabled"`
	AllowedFields []string `yaml:"allowed_fields"`
	Operators     []string `yaml:"operators"` // eq, ne, gt, lt, gte, lte, like, in
}

type ExportConfig struct {
	Enabled bool     `yaml:"enabled"`
	Formats []string `yaml:"formats"` // csv, xlsx, json, xml, pdf
	MaxRows int      `yaml:"max_rows"`
	Async   bool     `yaml:"async"` // 异步导出
}

type LoggingConfig struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`  // debug, info, warn, error
	Format  string `yaml:"format"` // json, text
	Output  string `yaml:"output"` // stdout, file, both
}

type ValidatorConfig struct {
	CustomRules   map[string]string `yaml:"custom_rules"`
	ErrorFormat   string            `yaml:"error_format"` // simple, detailed
	LocaleSupport bool              `yaml:"locale_support"`
}

// 新增：钩子配置
type HooksConfig struct {
	BeforeCreate []string `yaml:"before_create"`
	AfterCreate  []string `yaml:"after_create"`
	BeforeUpdate []string `yaml:"before_update"`
	AfterUpdate  []string `yaml:"after_update"`
	BeforeDelete []string `yaml:"before_delete"`
	AfterDelete  []string `yaml:"after_delete"`
	BeforeFind   []string `yaml:"before_find"`
	AfterFind    []string `yaml:"after_find"`
}

// 新增：关联关系配置
type RelationsConfig struct {
	Enabled   bool             `yaml:"enabled"`
	Relations []RelationConfig `yaml:"relations"`
}

type RelationConfig struct {
	Model      string `yaml:"model"`
	Type       string `yaml:"type"` // hasOne, hasMany, belongsTo, manyToMany
	ForeignKey string `yaml:"foreign_key"`
	References string `yaml:"references"`
	JoinTable  string `yaml:"join_table"` // for manyToMany
	Preload    bool   `yaml:"preload"`
}

// 新增：认证配置
type AuthConfig struct {
	Enabled   bool                `yaml:"enabled"`
	Type      string              `yaml:"type"` // jwt, oauth2, basic, apikey
	JWT       JWTConfig           `yaml:"jwt"`
	Protected []string            `yaml:"protected"` // 需要认证的方法
	Roles     map[string][]string `yaml:"roles"`     // 角色权限配置
}

type JWTConfig struct {
	Secret       string `yaml:"secret"`
	ExpireHours  int    `yaml:"expire_hours"`
	RefreshToken bool   `yaml:"refresh_token"`
	Algorithm    string `yaml:"algorithm"` // HS256, RS256
}

type ModelInfo struct {
	Name         string
	Fields       []FieldInfo
	PkgName      string
	LowerName    string
	PluralName   string
	SnakeName    string // 新增
	KebabName    string // 新增
	Config       *Config
	HasTimeField bool
	Relations    []RelationConfig // 新增
	Indexes      []IndexConfig    // 新增
	TableName    string           // 新增
}

type FieldInfo struct {
	Name        string
	Type        string
	JsonTag     string
	GormTag     string
	ValidateTag string // 新增：验证标签
	BindingTag  string // 新增：绑定标签
	IsID        bool
	Required    bool
	Searchable  bool
	Filterable  bool   // 新增：可过滤
	Sortable    bool   // 新增：可排序
	Unique      bool   // 新增：唯一
	Index       bool   // 新增：索引
	Default     string // 新增：默认值
	Comment     string // 新增：注释
	Encrypted   bool   // 新增：是否加密
	Sensitive   bool   // 新增：敏感字段（日志中不显示）
}

func DefaultConfig() *Config {
	return &Config{
		Model: ModelConfig{
			SoftDelete: false,
			Timestamps: true,
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
			RateLimit: RateLimitConfig{
				Enabled: false,
				Rate:    "100/1m",
				Burst:   10,
			},
			CORS: CORSConfig{
				Enabled: false,
			},
			Swagger: SwaggerConfig{
				Enabled:     false,
				Title:       "API Documentation",
				Description: "Auto-generated API Documentation",
				Version:     "1.0.0",
			},
		},
		Database: DatabaseConfig{
			Type:        "mysql",
			TimeFields:  true,
			Transaction: true,
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
				"handler":    "handlers",
				"route":      "routes",
				"service":    "services",
				"repository": "repositories",
				"dto":        "dto",
				"middleware": "middleware",
			},
			Overwrite: false,
			Format: FormatConfig{
				Enabled:   true,
				GoFmt:     true,
				GoImports: true,
			},
		},
		Features: FeaturesConfig{
			Validation: true,
			Cache: CacheConfig{
				Enabled:   false,
				Type:      "memory",
				TTL:       "5m",
				KeyPrefix: "api:",
				Strategy:  "cache-aside",
			},
			Search: false,
			Sort:   false,
			Filter: FilterConfig{
				Enabled:   false,
				Operators: []string{"eq", "ne", "gt", "lt", "like"},
			},
			Export: ExportConfig{
				Enabled: false,
				Formats: []string{"csv", "json"},
				MaxRows: 10000,
			},
			Service:    true,
			Repository: true,
			DTO:        true,
			Logging: LoggingConfig{
				Enabled: true,
				Level:   "info",
				Format:  "json",
				Output:  "stdout",
			},
			Metrics:  false,
			Recovery: true,
			Test:     false,
			Mock:     false,
		},
		Relations: RelationsConfig{
			Enabled: false,
		},
		Auth: AuthConfig{
			Enabled: false,
			Type:    "jwt",
			JWT: JWTConfig{
				ExpireHours: 24,
				Algorithm:   "HS256",
			},
		},
	}
}

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
