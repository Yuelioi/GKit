package i18n

// NewBuiltinRegistry 返回官方内置翻译
func NewBuiltinRegistry() *Registry {
	r := NewRegistry()

	r.RegisterBatch(EN, map[Key]string{
		"error.success":        "Success",
		"error.invalid_params": "Invalid parameters",
		"error.unauthorized":   "Unauthorized",
		"error.forbidden":      "Forbidden",
		"error.not_found":      "Not found",
		"error.internal":       "Internal server error",
	})

	r.RegisterBatch(ZH, map[Key]string{
		"error.success":        "成功",
		"error.invalid_params": "参数错误",
		"error.unauthorized":   "未授权",
		"error.forbidden":      "禁止访问",
		"error.not_found":      "资源不存在",
		"error.internal":       "服务器内部错误",
	})

	return r
}
