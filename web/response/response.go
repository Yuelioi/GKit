package response

import (
	"fmt"
	"time"
)

// Response 统一响应结构
type Response struct {
	Code      int    `json:"code"`                 // 统一状态码 (0为成功，非0为各种业务/系统错误)
	Message   string `json:"message"`              // 响应消息，对用户的友好描述
	Data      any    `json:"data,omitempty"`       // 响应数据
	RequestID string `json:"request_id,omitempty"` // 请求追踪ID
	Timestamp int64  `json:"timestamp"`            // 时间戳（毫秒级）
	Error     string `json:"-"`                    // 详细错误信息（内部使用，不返回给客户端）
}

// CodeInfo 状态码信息
type CodeInfo struct {
	Code       int    // 状态码
	Message    string // 默认错误消息
	HTTPStatus int    // 对应的HTTP状态码
}

// 预定义的状态码
var (
	// 成功
	Success = CodeInfo{0, "操作成功", 200}

	// 1XXXX: 客户端请求/参数错误
	InvalidParam     = CodeInfo{10001, "请求参数错误或缺失", 400}
	ResourceExist    = CodeInfo{10002, "资源已存在，请勿重复创建", 400}
	MissingParam     = CodeInfo{10003, "缺少必填参数", 400}
	InvalidFormat    = CodeInfo{10004, "参数格式不正确", 400}
	TooManyRequest   = CodeInfo{10005, "请求过于频繁，请稍后重试", 429}
	RequestTimeout   = CodeInfo{10006, "请求超时，请重试", 408}
	FileTooLarge     = CodeInfo{10007, "文件大小超出限制", 400}
	InvalidFileType  = CodeInfo{10008, "不支持的文件类型", 400}
	InvalidEmail     = CodeInfo{10009, "邮箱格式不正确", 400}
	InvalidPhone     = CodeInfo{10010, "手机号格式不正确", 400}
	InvalidIDCard    = CodeInfo{10011, "身份证号格式不正确", 400}
	InvalidURL       = CodeInfo{10012, "URL格式不正确", 400}
	InvalidDate      = CodeInfo{10013, "日期格式不正确", 400}
	OutOfRange       = CodeInfo{10014, "参数值超出允许范围", 400}
	Duplicate        = CodeInfo{10015, "请勿重复提交", 400}
	OperationFailed  = CodeInfo{10016, "操作失败，请重试", 400}
	InvalidOperation = CodeInfo{10017, "当前操作无效", 400}
	InvalidState     = CodeInfo{10018, "当前状态不允许此操作", 400}
	VersionConflict  = CodeInfo{10019, "数据版本冲突，请刷新后重试", 409}

	// 2XXXX: 业务逻辑错误
	BusinessError       = CodeInfo{20001, "业务处理失败", 400}
	AccountDisabled     = CodeInfo{20002, "账号已被禁用，请联系管理员", 403}
	AccountNotExist     = CodeInfo{20003, "账号不存在", 404}
	PasswordError       = CodeInfo{20004, "密码错误", 400}
	PasswordWeak        = CodeInfo{20005, "密码强度不够，请使用更复杂的密码", 400}
	VerifyCodeError     = CodeInfo{20006, "验证码错误", 400}
	VerifyCodeExpired   = CodeInfo{20007, "验证码已过期，请重新获取", 400}
	InsufficientBalance = CodeInfo{20008, "余额不足", 400}
	InsufficientStock   = CodeInfo{20009, "库存不足", 400}
	OrderNotExist       = CodeInfo{20010, "订单不存在", 404}
	OrderPaid           = CodeInfo{20011, "订单已支付，请勿重复支付", 400}
	OrderCancelled      = CodeInfo{20012, "订单已取消", 400}
	PaymentFailed       = CodeInfo{20013, "支付失败，请重试", 400}
	RefundFailed        = CodeInfo{20014, "退款失败，请联系客服", 400}
	AccountLocked       = CodeInfo{20015, "账号已锁定，请稍后再试", 403}
	TooManyAttempts     = CodeInfo{20016, "尝试次数过多，请稍后再试", 429}
	OperationNotAllowed = CodeInfo{20017, "当前操作不被允许", 403}
	QuotaExceeded       = CodeInfo{20018, "使用配额已用尽", 429}
	ServiceExpired      = CodeInfo{20019, "服务已过期，请续费", 403}
	NotInServiceTime    = CodeInfo{20020, "当前不在服务时间", 400}

	// 3XXXX: 第三方服务错误
	SMSError        = CodeInfo{30001, "短信发送失败，请稍后重试", 500}
	EmailError      = CodeInfo{30002, "邮件发送失败，请稍后重试", 500}
	PaymentError    = CodeInfo{30003, "支付服务异常，请稍后重试", 500}
	StorageError    = CodeInfo{30004, "存储服务异常", 500}
	ThirdPartyError = CodeInfo{30005, "第三方服务异常，请稍后重试", 500}
	UploadFailed    = CodeInfo{30006, "文件上传失败", 500}
	DownloadFailed  = CodeInfo{30007, "文件下载失败", 500}

	// 4XXXX: 鉴权/权限相关错误
	Unauthorized   = CodeInfo{40001, "请先登录或登录信息已过期", 401}
	TokenExpired   = CodeInfo{40002, "登录已过期，请重新登录", 401}
	Forbidden      = CodeInfo{40003, "权限不足，无法执行该操作", 403}
	NotFound       = CodeInfo{40004, "请求的资源不存在", 404}
	InvalidToken   = CodeInfo{40005, "无效的访问令牌", 401}
	TokenRevoked   = CodeInfo{40006, "访问令牌已被吊销", 401}
	SignatureError = CodeInfo{40007, "签名验证失败", 401}
	IPBlocked      = CodeInfo{40008, "您的IP已被封禁", 403}
	AccountBanned  = CodeInfo{40009, "账号已被封禁，请联系管理员", 403}
	RoleNotAllowed = CodeInfo{40010, "当前角色无权执行此操作", 403}

	// 5XXXX: 数据相关错误
	DataNotFound     = CodeInfo{50001, "数据不存在", 404}
	DataExists       = CodeInfo{50002, "数据已存在", 409}
	DataInvalid      = CodeInfo{50003, "数据格式无效", 400}
	DataExpired      = CodeInfo{50004, "数据已过期", 410}
	DataCorrupted    = CodeInfo{50005, "数据已损坏", 500}
	DataInconsistent = CodeInfo{50006, "数据不一致，请刷新后重试", 409}

	// 9XXXX: 服务器系统/内部错误
	InternalError      = CodeInfo{90001, "服务器开小差了，请稍后重试", 500}
	ServiceBusy        = CodeInfo{90002, "系统繁忙，请稍后再试", 503}
	DatabaseError      = CodeInfo{90003, "数据库操作失败", 500}
	CacheError         = CodeInfo{90004, "缓存操作失败", 500}
	RemoteCallError    = CodeInfo{90005, "远程服务调用失败", 500}
	ServiceUnavailable = CodeInfo{90006, "服务暂时不可用", 503}
	NetworkError       = CodeInfo{90007, "网络连接失败", 500}
	ConfigError        = CodeInfo{90008, "系统配置错误", 500}
	DependencyError    = CodeInfo{90009, "依赖服务异常", 500}
	MaintenanceMode    = CodeInfo{90010, "系统维护中，请稍后访问", 503}
	ResourceExhausted  = CodeInfo{90011, "系统资源不足", 503}
	Deadlock           = CodeInfo{90012, "系统繁忙，请重试", 500}
	TransactionError   = CodeInfo{90013, "事务处理失败", 500}
)

// codeMap 状态码映射表（用于通过code快速查找CodeInfo）
var codeMap = map[int]CodeInfo{
	Success.Code: Success,

	// 1XXXX
	InvalidParam.Code:     InvalidParam,
	ResourceExist.Code:    ResourceExist,
	MissingParam.Code:     MissingParam,
	InvalidFormat.Code:    InvalidFormat,
	TooManyRequest.Code:   TooManyRequest,
	RequestTimeout.Code:   RequestTimeout,
	FileTooLarge.Code:     FileTooLarge,
	InvalidFileType.Code:  InvalidFileType,
	InvalidEmail.Code:     InvalidEmail,
	InvalidPhone.Code:     InvalidPhone,
	InvalidIDCard.Code:    InvalidIDCard,
	InvalidURL.Code:       InvalidURL,
	InvalidDate.Code:      InvalidDate,
	OutOfRange.Code:       OutOfRange,
	Duplicate.Code:        Duplicate,
	OperationFailed.Code:  OperationFailed,
	InvalidOperation.Code: InvalidOperation,
	InvalidState.Code:     InvalidState,
	VersionConflict.Code:  VersionConflict,

	// 2XXXX
	BusinessError.Code:       BusinessError,
	AccountDisabled.Code:     AccountDisabled,
	AccountNotExist.Code:     AccountNotExist,
	PasswordError.Code:       PasswordError,
	PasswordWeak.Code:        PasswordWeak,
	VerifyCodeError.Code:     VerifyCodeError,
	VerifyCodeExpired.Code:   VerifyCodeExpired,
	InsufficientBalance.Code: InsufficientBalance,
	InsufficientStock.Code:   InsufficientStock,
	OrderNotExist.Code:       OrderNotExist,
	OrderPaid.Code:           OrderPaid,
	OrderCancelled.Code:      OrderCancelled,
	PaymentFailed.Code:       PaymentFailed,
	RefundFailed.Code:        RefundFailed,
	AccountLocked.Code:       AccountLocked,
	TooManyAttempts.Code:     TooManyAttempts,
	OperationNotAllowed.Code: OperationNotAllowed,
	QuotaExceeded.Code:       QuotaExceeded,
	ServiceExpired.Code:      ServiceExpired,
	NotInServiceTime.Code:    NotInServiceTime,

	// 3XXXX
	SMSError.Code:        SMSError,
	EmailError.Code:      EmailError,
	PaymentError.Code:    PaymentError,
	StorageError.Code:    StorageError,
	ThirdPartyError.Code: ThirdPartyError,
	UploadFailed.Code:    UploadFailed,
	DownloadFailed.Code:  DownloadFailed,

	// 4XXXX
	Unauthorized.Code:   Unauthorized,
	TokenExpired.Code:   TokenExpired,
	Forbidden.Code:      Forbidden,
	NotFound.Code:       NotFound,
	InvalidToken.Code:   InvalidToken,
	TokenRevoked.Code:   TokenRevoked,
	SignatureError.Code: SignatureError,
	IPBlocked.Code:      IPBlocked,
	AccountBanned.Code:  AccountBanned,
	RoleNotAllowed.Code: RoleNotAllowed,

	// 5XXXX
	DataNotFound.Code:     DataNotFound,
	DataExists.Code:       DataExists,
	DataInvalid.Code:      DataInvalid,
	DataExpired.Code:      DataExpired,
	DataCorrupted.Code:    DataCorrupted,
	DataInconsistent.Code: DataInconsistent,

	// 9XXXX
	InternalError.Code:      InternalError,
	ServiceBusy.Code:        ServiceBusy,
	DatabaseError.Code:      DatabaseError,
	CacheError.Code:         CacheError,
	RemoteCallError.Code:    RemoteCallError,
	ServiceUnavailable.Code: ServiceUnavailable,
	NetworkError.Code:       NetworkError,
	ConfigError.Code:        ConfigError,
	DependencyError.Code:    DependencyError,
	MaintenanceMode.Code:    MaintenanceMode,
	ResourceExhausted.Code:  ResourceExhausted,
	Deadlock.Code:           Deadlock,
	TransactionError.Code:   TransactionError,
}

// ---------- 内部构造函数 ----------

// newResponse 基础构造函数
func newResponse(codeInfo CodeInfo, message string, data any) *Response {
	// 如果没有提供自定义消息，使用默认消息
	if message == "" {
		message = codeInfo.Message
	}

	return &Response{
		Code:      codeInfo.Code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

// ---------- 通用响应构造函数 ----------

// OK 成功响应（使用预定义的Success）
func OK(data any) *Response {
	return newResponse(Success, "", data)
}

// OKWithMsg 成功响应（自定义消息）
func OKWithMsg(message string, data any) *Response {
	return newResponse(Success, message, data)
}

// OKWithNoContent 成功响应（无数据）
func OKWithNoContent() *Response {
	return newResponse(Success, "操作成功", nil)
}

// Fail 错误响应（使用预定义的CodeInfo）
func Fail(codeInfo CodeInfo) *Response {
	return newResponse(codeInfo, "", nil)
}

// FailWithMsg 错误响应（自定义消息）
func FailWithMsg(codeInfo CodeInfo, message string) *Response {
	return newResponse(codeInfo, message, nil)
}
func FailWithData(codeInfo CodeInfo, data any) *Response {
	return newResponse(codeInfo, "", data)
}

// FailWithMsgAndData 错误响应（自定义消息和数据）
func FailWithMsgAndData(codeInfo CodeInfo, message string, data any) *Response {
	return newResponse(codeInfo, message, data)
}

// FailWithError 错误响应（附带error对象，用于内部日志）
func FailWithError(codeInfo CodeInfo, err error) *Response {
	resp := newResponse(codeInfo, "", nil)
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}

// ---------- 便捷方法 ----------

// NotFoundWithResource 资源不存在（可指定资源名称）
func NotFoundWithResource(resourceName string) *Response {
	msg := NotFound.Message
	if resourceName != "" {
		msg = fmt.Sprintf("%s不存在", resourceName)
	}
	return newResponse(NotFound, msg, nil)
}

// InvalidParamWithField 参数错误（指定字段名）
func InvalidParamWithField(fieldName string) *Response {
	msg := InvalidParam.Message
	if fieldName != "" {
		msg = fmt.Sprintf("参数 %s 错误", fieldName)
	}
	return newResponse(InvalidParam, msg, nil)
}

// MissingParamWithField 缺少参数（指定字段名）
func MissingParamWithField(fieldName string) *Response {
	msg := MissingParam.Message
	if fieldName != "" {
		msg = fmt.Sprintf("缺少必填参数: %s", fieldName)
	}
	return newResponse(MissingParam, msg, nil)
}

// ---------- 链式调用方法 ----------

// WithRequestID 添加请求ID
func (r *Response) WithRequestID(requestID string) *Response {
	r.RequestID = requestID
	return r
}

// WithTimestamp 更新为当前时间戳
func (r *Response) WithTimestamp() *Response {
	r.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	return r
}

// WithError 添加详细错误信息（内部使用）
func (r *Response) WithError(err error) *Response {
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

// WithMessage 覆盖消息
func (r *Response) WithMessage(message string) *Response {
	if message != "" {
		r.Message = message
	}
	return r
}

// WithData 设置数据
func (r *Response) WithData(data any) *Response {
	r.Data = data
	return r
}

// ---------- 工具方法 ----------

// HTTPStatus 获取对应的HTTP状态码
func (r *Response) HTTPStatus() int {
	if codeInfo, exists := codeMap[r.Code]; exists {
		return codeInfo.HTTPStatus
	}
	// 未知code，默认返回500
	return 500
}

// IsSuccess 判断是否成功
func (r *Response) IsSuccess() bool {
	return r.Code == Success.Code
}
