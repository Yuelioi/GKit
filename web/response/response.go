package response

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Version   string      `json:"version"`

	httpStatus int
}

func (r *Response) Status() int {
	if r.httpStatus > 0 {
		return r.httpStatus
	}
	return 200
}
