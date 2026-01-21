package errorx

func GetCause(err error) error {
	if c, ok := err.(Code); ok {
		return c.Cause()
	}
	return err
}

func Equal(err1, err2 error) bool {
	c1, ok1 := err1.(Code)
	c2, ok2 := err2.(Code)

	if !ok1 || !ok2 {
		return false
	}

	return c1.Code() == c2.Code()
}

func IsRetriable(err error) bool {
	c, ok := err.(Code)
	if !ok {
		return false
	}
	return c.IsRetriable()
}
