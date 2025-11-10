package response_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Yuelioi/gkit/web/response"
)

// mockJSON 模拟 Gin 的 c.JSON 函数
func mockJSON(status int, resp *response.Response) {
	fmt.Printf("HTTP Status: %d\n", status)
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
	fmt.Println("--------------------------------------------------")
}

func TestResponseExamples(t *testing.T) {
	// 成功响应
	data := map[string]interface{}{
		"id":   123,
		"name": "Alice",
	}
	mockJSON(200, response.Success(data).WithRequestID("req-001"))

	// 成功但无返回数据
	mockJSON(200, response.NoContent().WithRequestID("req-002"))

	// 参数错误
	mockJSON(400, response.InvalidParam("用户名不能为空").WithRequestID("req-003"))

	// 未授权
	mockJSON(401, response.Unauthorized("").WithRequestID("req-004"))

	// 资源不存在
	mockJSON(404, response.NotFound("用户").WithRequestID("req-005"))

	// 内部错误
	err := errors.New("数据库连接失败")
	mockJSON(500, response.InternalError(err).WithRequestID("req-006"))

	// 自定义错误
	mockJSON(400, response.Error(response.CodeResourceExist, "订单已存在").WithRequestID("req-007"))

	// 带时间戳更新
	mockJSON(200, response.Success(map[string]string{"status": "ok"}).WithTimestamp())

	// 演示多次打印时间戳变化
	time.Sleep(time.Millisecond * 10)
	mockJSON(200, response.Success("再次测试时间戳").WithTimestamp())
}
