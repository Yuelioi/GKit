package errorx

import "fmt"

// ============ 错误码注册表（核心机制） ============

// CodeRegistry 错误码注册表
// 每个应用可以注册自己的错误码，支持动态查询和验证
type CodeRegistry struct {
	codes map[int]CodeSpec
}

// CodeSpec 错误码规范
type CodeSpec struct {
	Code       int
	Message    string
	HttpStatus int
	Desc       string // 描述，用于文档生成
	Retriable  bool   // 是否可重试
}

var globalRegistry = &CodeRegistry{
	codes: make(map[int]CodeSpec),
}

// Register 注册错误码
func Register(spec CodeSpec) {
	if _, exists := globalRegistry.codes[spec.Code]; exists {
		panic(fmt.Sprintf("error code %d already registered", spec.Code))
	}
	globalRegistry.codes[spec.Code] = spec
}

// RegisterBatch 批量注册错误码
func RegisterBatch(specs []CodeSpec) {
	for _, spec := range specs {
		Register(spec)
	}
}

// GetSpec 获取错误码规范（用于文档生成、验证等）
func GetSpec(code int) (CodeSpec, bool) {
	spec, ok := globalRegistry.codes[code]
	return spec, ok
}

// ListSpecs 列出所有错误码（用于文档生成）
func ListSpecs() []CodeSpec {
	specs := make([]CodeSpec, 0, len(globalRegistry.codes))
	for _, spec := range globalRegistry.codes {
		specs = append(specs, spec)
	}
	return specs
}
