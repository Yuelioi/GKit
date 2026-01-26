// response/response.go
package response

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`

	// 内部字段，不序列化到 JSON
	httpStatus int `json:"-"`
}
