package errorx

type CodeSpec struct {
	Code       int
	MessageKey string
	HttpStatus int
	Desc       string
	Version    string
	Retriable  bool
}
