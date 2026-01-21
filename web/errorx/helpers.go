package errorx

func Is(err error, target Code) bool {
	e, ok1 := err.(Code)
	t, ok2 := target.(Code)
	return ok1 && ok2 && e.Code() == t.Code()
}

func Cause(err error) error {
	if c, ok := err.(Code); ok {
		return c.Cause()
	}
	return err
}
