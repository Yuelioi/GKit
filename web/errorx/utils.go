package errorx

// ============ 工具函数 ============

// GetCause 获取原始错误
func GetCause(err error) error {
	if e, ok := err.(*codeErr); ok && e.cause != nil {
		return e.cause
	}
	return err
}

// Equal 比较错误码
func Equal(err1, err2 error) bool {
	c1, ok1 := err1.(Code)
	c2, ok2 := err2.(Code)

	if !ok1 || !ok2 {
		return false
	}

	return c1.Code() == c2.Code()
}

// Is 判断是否为指定错误
func Is(err error, target Code) bool {
	return Equal(err, target)
}

// IsRetriable 判断错误是否可重试
// 应用启动时注册错误码时指定 Retriable 字段
func IsRetriable(err error) bool {
	c, ok := err.(Code)
	if !ok {
		return false
	}

	spec, ok := globalRegistry.codes[c.Code()]
	if !ok {
		// 未注册的错误默认不重试
		return false
	}

	return spec.Retriable
}
