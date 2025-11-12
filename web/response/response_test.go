package response_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Yuelioi/gkit/web/response"
)

// mockJSON 模拟 Gin 的 c.JSON 函数
func mockJSON(status int, resp *response.Response) {
	fmt.Printf("HTTP Status: %d\n", status)
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
	fmt.Println("--------------------------------------------------")
}

func TestSuccessResponse(t *testing.T) {
	fmt.Println("========== 成功响应示例 ==========")

	// 1. 返回数据
	userData := map[string]interface{}{
		"id":       123,
		"username": "alice",
		"email":    "alice@example.com",
	}
	resp := response.OK(userData).WithRequestID("req-001")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 返回列表数据
	listData := []map[string]interface{}{
		{"id": 1, "name": "Item 1"},
		{"id": 2, "name": "Item 2"},
	}
	resp = response.OK(listData).WithRequestID("req-002")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 成功但无返回数据（如删除操作）
	resp = response.OKWithNoContent().WithRequestID("req-003")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 自定义成功消息
	resp = response.OKWithMsg("注册成功，请查收验证邮件", map[string]string{
		"email": "alice@example.com",
	}).WithRequestID("req-004")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestClientErrorResponse(t *testing.T) {
	fmt.Println("\n========== 客户端错误响应示例 ==========")

	// 1. 参数错误 - 使用默认消息
	resp := response.Fail(response.InvalidParam).WithRequestID("req-101")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 参数错误 - 指定字段名
	resp = response.InvalidParamWithField("email").WithRequestID("req-102")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 参数错误 - 自定义消息
	resp = response.FailWithMsg(response.InvalidParam, "邮箱格式不正确，请输入有效的邮箱地址").
		WithRequestID("req-103")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 缺少必填参数
	resp = response.MissingParamWithField("password").WithRequestID("req-104")
	mockJSON(resp.HTTPStatus(), resp)

	// 5. 文件相关错误
	resp = response.Fail(response.FileTooLarge).
		WithMessage("文件大小不能超过10MB").
		WithRequestID("req-105")
	mockJSON(resp.HTTPStatus(), resp)

	// 6. 请求频繁
	resp = response.Fail(response.TooManyRequest).
		WithMessage("您的操作过于频繁，请在60秒后重试").
		WithRequestID("req-106")
	mockJSON(resp.HTTPStatus(), resp)

	// 7. 重复提交
	resp = response.Fail(response.Duplicate).WithRequestID("req-107")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestAuthErrorResponse(t *testing.T) {
	fmt.Println("\n========== 鉴权/权限错误响应示例 ==========")

	// 1. 未登录
	resp := response.Fail(response.Unauthorized).WithRequestID("req-201")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. Token过期
	resp = response.Fail(response.TokenExpired).
		WithMessage("登录已过期，请重新登录以继续操作").
		WithRequestID("req-202")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 权限不足
	resp = response.FailWithMsg(response.Forbidden, "您没有删除该资源的权限").
		WithRequestID("req-203")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 资源不存在
	resp = response.NotFoundWithResource("订单").WithRequestID("req-204")
	mockJSON(resp.HTTPStatus(), resp)

	// 5. 账号被封禁
	resp = response.Fail(response.AccountBanned).
		WithMessage("您的账号因违规操作已被封禁，如有疑问请联系客服").
		WithRequestID("req-205")
	mockJSON(resp.HTTPStatus(), resp)

	// 6. IP被封禁
	resp = response.Fail(response.IPBlocked).WithRequestID("req-206")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestBusinessErrorResponse(t *testing.T) {
	fmt.Println("\n========== 业务逻辑错误响应示例 ==========")

	// 1. 账号不存在
	resp := response.Fail(response.AccountNotExist).WithRequestID("req-301")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 密码错误（带尝试次数提示）
	resp = response.FailWithMsg(response.PasswordError, "密码错误，您还有2次尝试机会").
		WithRequestID("req-302")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 验证码错误
	resp = response.Fail(response.VerifyCodeError).WithRequestID("req-303")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 余额不足（带数据）
	resp = response.FailWithData(response.InsufficientBalance, map[string]interface{}{
		"balance":  50.00,
		"required": 100.00,
	}).WithRequestID("req-304")
	mockJSON(resp.HTTPStatus(), resp)

	// 5. 库存不足
	resp = response.FailWithMsg(response.InsufficientStock, "商品库存不足，仅剩5件").
		WithRequestID("req-305")
	mockJSON(resp.HTTPStatus(), resp)

	// 6. 订单已支付
	resp = response.Fail(response.OrderPaid).WithRequestID("req-306")
	mockJSON(resp.HTTPStatus(), resp)

	// 7. 配额用尽
	resp = response.FailWithMsg(response.QuotaExceeded, "今日API调用次数已用尽，请明天再试或升级套餐").
		WithRequestID("req-307")
	mockJSON(resp.HTTPStatus(), resp)

	// 8. 账号锁定
	resp = response.Fail(response.AccountLocked).
		WithMessage("密码错误次数过多，账号已锁定30分钟").
		WithRequestID("req-308")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestServerErrorResponse(t *testing.T) {
	fmt.Println("\n========== 服务器错误响应示例 ==========")

	// 1. 内部错误（不暴露具体错误）
	resp := response.Fail(response.InternalError).WithRequestID("req-401")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 内部错误（记录详细错误日志）
	err := errors.New("database connection timeout: connection refused")
	resp = response.FailWithError(response.DatabaseError, err).WithRequestID("req-402")
	mockJSON(resp.HTTPStatus(), resp)
	// 注意：resp.Error 字段不会返回给客户端（json:"-"），仅用于日志
	if resp.Error != "" {
		fmt.Printf("内部日志: %s\n", resp.Error)
		fmt.Println("--------------------------------------------------")
	}

	// 3. 服务繁忙
	resp = response.Fail(response.ServiceBusy).WithRequestID("req-403")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 远程调用失败
	resp = response.FailWithMsg(response.RemoteCallError, "支付服务暂时不可用，请稍后重试").
		WithRequestID("req-404")
	mockJSON(resp.HTTPStatus(), resp)

	// 5. 系统维护
	resp = response.Fail(response.MaintenanceMode).
		WithMessage("系统维护中(预计30分钟)，请稍后访问").
		WithRequestID("req-405")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestThirdPartyErrorResponse(t *testing.T) {
	fmt.Println("\n========== 第三方服务错误响应示例 ==========")

	// 1. 短信发送失败
	resp := response.Fail(response.SMSError).
		WithMessage("验证码发送失败，请检查手机号或稍后重试").
		WithRequestID("req-501")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 邮件发送失败
	resp = response.Fail(response.EmailError).WithRequestID("req-502")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 文件上传失败
	resp = response.FailWithMsg(response.UploadFailed, "图片上传失败，请重试").
		WithRequestID("req-503")
	mockJSON(resp.HTTPStatus(), resp)

	// 4. 支付服务异常
	resp = response.Fail(response.PaymentError).WithRequestID("req-504")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestChainedCalls(t *testing.T) {
	fmt.Println("\n========== 链式调用示例 ==========")

	// 1. 多个链式调用
	resp := response.Fail(response.InvalidParam).
		WithMessage("用户名长度必须在3-20个字符之间").
		WithRequestID("req-601").
		WithTimestamp()
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 成功响应的链式调用
	userData := map[string]interface{}{
		"id":       456,
		"username": "bob",
	}
	resp = response.OK(userData).
		WithMessage("登录成功").
		WithRequestID("req-602").
		WithTimestamp()
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 错误响应添加额外数据
	resp = response.Fail(response.PasswordWeak).
		WithMessage("密码强度不够").
		WithData(map[string]interface{}{
			"requirements": []string{
				"至少8个字符",
				"包含大小写字母",
				"包含数字",
				"包含特殊字符",
			},
		}).
		WithRequestID("req-603")
	mockJSON(resp.HTTPStatus(), resp)
}

func TestComplexScenarios(t *testing.T) {
	fmt.Println("\n========== 复杂场景示例 ==========")

	// 1. 用户注册场景
	fmt.Println("场景1: 用户注册")
	resp := response.FailWithMsgAndData(
		response.InvalidParam,
		"注册失败，请检查以下信息",
		map[string]interface{}{
			"errors": []map[string]string{
				{"field": "email", "message": "邮箱格式不正确"},
				{"field": "password", "message": "密码强度不够"},
			},
		},
	).WithRequestID("req-701")
	mockJSON(resp.HTTPStatus(), resp)

	// 2. 订单创建场景
	fmt.Println("场景2: 订单创建失败")
	resp = response.FailWithData(
		response.InsufficientStock,
		map[string]interface{}{
			"product_id": "P12345",
			"requested":  10,
			"available":  3,
			"message":    "该商品库存不足，当前仅剩3件",
		},
	).WithRequestID("req-702")
	mockJSON(resp.HTTPStatus(), resp)

	// 3. 批量操作场景
	fmt.Println("场景3: 批量删除部分成功")
	resp = response.OKWithMsg(
		"批量删除完成，部分项目删除失败",
		map[string]interface{}{
			"total":   10,
			"success": 7,
			"failed":  3,
			"failed_items": []map[string]interface{}{
				{"id": 5, "reason": "权限不足"},
				{"id": 8, "reason": "资源不存在"},
				{"id": 9, "reason": "资源被占用"},
			},
		},
	).WithRequestID("req-703")
	mockJSON(resp.HTTPStatus(), resp)
}

// TestGinIntegration 模拟 Gin 框架集成
func TestGinIntegration(t *testing.T) {
	fmt.Println("\n========== Gin 框架集成示例 ==========")

	// 模拟 Gin 的 Context
	type MockGinContext struct{}
	c := &MockGinContext{}

	// 模拟 c.JSON 方法
	mockGinJSON := func(c *MockGinContext, resp *response.Response) {
		mockJSON(resp.HTTPStatus(), resp)
	}

	// 示例1: 获取用户信息
	fmt.Println("GET /api/users/123")
	user := map[string]interface{}{
		"id":       123,
		"username": "alice",
		"email":    "alice@example.com",
		"role":     "admin",
	}
	mockGinJSON(c, response.OK(user).WithRequestID("req-801"))

	// 示例2: 创建用户（参数错误）
	fmt.Println("POST /api/users")
	mockGinJSON(c, response.InvalidParamWithField("email").WithRequestID("req-802"))

	// 示例3: 更新用户（权限不足）
	fmt.Println("PUT /api/users/456")
	mockGinJSON(c, response.Fail(response.Forbidden).
		WithMessage("只能修改自己的信息").
		WithRequestID("req-803"))

	// 示例4: 删除用户（成功）
	fmt.Println("DELETE /api/users/789")
	mockGinJSON(c, response.OKWithNoContent().WithRequestID("req-804"))

	// 示例5: 登录（密码错误）
	fmt.Println("POST /api/auth/login")
	mockGinJSON(c, response.FailWithMsg(response.PasswordError, "密码错误，您还有2次尝试机会").
		WithRequestID("req-805"))

	// 示例6: 数据库错误（内部错误）
	fmt.Println("GET /api/orders")
	dbErr := errors.New("sql: connection refused")
	mockGinJSON(c, response.FailWithError(response.DatabaseError, dbErr).
		WithRequestID("req-806"))
}
