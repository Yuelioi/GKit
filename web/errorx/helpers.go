package errorx

// ============ 错误判断 ============

func Is(err error, target Error) bool {
	e, ok := err.(Error)
	return ok && e.Code() == target.Code()
}

// Cause 递归获取最底层的原始错误
func Cause(err error) error {
	for {
		if e, ok := err.(Error); ok {
			if cause := e.Cause(); cause != nil {
				err = cause
				continue
			}
		}
		return err
	}
}

// IsRetriable 判断是否可重试
func IsRetriable(err error) bool {
	if e, ok := err.(Error); ok {
		return e.IsRetriable()
	}
	return false
}
