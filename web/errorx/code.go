package errorx

type Code interface {
	error
	Code() int
	MessageKey() string
	HttpStatus() int
	Version() string
	IsRetriable() bool
	Cause() error
}

type codeErr struct {
	spec  CodeSpec
	cause error
}

func (e *codeErr) Error() string      { return e.spec.MessageKey }
func (e *codeErr) Code() int          { return e.spec.Code }
func (e *codeErr) MessageKey() string { return e.spec.MessageKey }
func (e *codeErr) HttpStatus() int    { return e.spec.HttpStatus }
func (e *codeErr) Version() string    { return e.spec.Version }
func (e *codeErr) IsRetriable() bool  { return e.spec.Retriable }
func (e *codeErr) Cause() error       { return e.cause }

func New(spec CodeSpec) Code {
	return &codeErr{spec: spec}
}

func Wrap(c Code, cause error) Code {
	return &codeErr{
		spec:  GetSpecMust(c.Code()),
		cause: cause,
	}
}
