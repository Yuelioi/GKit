package response_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Yuelioi/gkit/web/response"
)

// ========== 测试数据结构 ==========

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// ========== 辅助函数 ==========

// printResponse 格式化输出响应
func printResponse(title string, resp *response.Response) {
	fmt.Printf("\n【%s】\n", title)
	fmt.Printf("HTTP Status: %d\n", resp.Status())
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
	fmt.Println("--------------------------------------------------")
}

// ========== 成功响应测试 ==========

func TestSuccessResponses(t *testing.T) {
	fmt.Println("\n========== 成功响应示例 ==========")

	// 1. 返回对象数据
	user := User{
		ID:       1,
		Username: "张三",
		Email:    "zhangsan@example.com",
		Age:      25,
	}
	printResponse("成功返回用户数据", response.Data(user))

	// 2. 返回列表数据
	users := []User{
		{ID: 1, Username: "张三", Email: "zhangsan@example.com", Age: 25},
		{ID: 2, Username: "李四", Email: "lisi@example.com", Age: 30},
		{ID: 3, Username: "王五", Email: "wangwu@example.com", Age: 28},
	}
	printResponse("成功返回用户列表", response.Data(users))

	// 3. 分页数据
	printResponse("分页查询结果",
		response.Page(users, 100, 1, 10))

	// 4. 创建资源成功 (201)
	newUser := User{ID: 4, Username: "赵六", Email: "zhaoliu@example.com", Age: 22}
	printResponse("创建用户成功", response.Created(newUser))

	// 5. 删除成功 (204 无内容)
	printResponse("删除成功", response.NoContent())

	// 6. 自定义消息
	printResponse("更新成功(自定义消息)",
		response.Data(user).WithMessage("用户信息更新成功"))

	// 7. 带请求ID
	printResponse("带请求追踪ID",
		response.Data(user).WithRequestID("req-abc123xyz"))

	// 8. 链式调用
	printResponse("链式调用示例",
		response.Data(user).
			WithMessage("查询成功").
			WithRequestID("req-chain-001"))

	// 9. 空数据成功
	printResponse("成功但无数据", response.OK())

	// 10. 返回简单类型
	printResponse("返回字符串", response.Data("操作执行成功"))
	printResponse("返回数字", response.Data(42))
	printResponse("返回布尔值", response.Data(true))
	printResponse("返回Map", response.Data(map[string]any{
		"total_users":  1000,
		"active_users": 856,
		"growth_rate":  12.5,
	}))
}

// ========== 错误响应测试 ==========

func TestErrorResponses(t *testing.T) {
	fmt.Println("\n========== 错误响应示例 ==========")

	// 1. 通用错误
	printResponse("通用错误", response.Error("操作失败"))

	// 2. 参数错误 (400)
	printResponse("参数错误", response.BadRequest("用户名不能为空"))
	printResponse("参数错误(默认消息)", response.BadRequest(""))

	// 3. 未授权 (401)
	printResponse("未授权", response.Unauthorized("请先登录"))
	printResponse("Token过期", response.Unauthorized("token已过期，请重新登录"))

	// 4. 禁止访问 (403)
	printResponse("权限不足", response.Forbidden("您没有权限访问该资源"))
	printResponse("禁止访问(默认消息)", response.Forbidden(""))

	// 5. 资源不存在 (404)
	printResponse("资源不存在", response.NotFound("用户不存在"))
	printResponse("接口不存在", response.NotFound("请求的接口不存在"))

	// 6. 服务器错误 (500)
	printResponse("内部错误", response.InternalError("数据库连接失败"))
	printResponse("内部错误(默认消息)", response.InternalError(""))

	// 7. 服务不可用 (503)
	printResponse("服务不可用", response.ServiceUnavailable("系统维护中"))
	printResponse("服务不可用(默认)", response.ServiceUnavailable(""))
}

// ========== 业务错误测试 ==========

func TestBusinessErrors(t *testing.T) {
	fmt.Println("\n========== 业务错误示例 ==========")

	// 1. 通用业务错误
	printResponse("业务错误", response.BusinessError("余额不足，无法完成支付"))

	// 2. 记录已存在
	printResponse("记录已存在", response.RecordExists("用户名已被注册"))
	printResponse("记录已存在(默认)", response.RecordExists(""))

	// 3. 记录不存在
	printResponse("记录不存在", response.RecordNotFound("订单不存在"))
	printResponse("记录不存在(默认)", response.RecordNotFound(""))

	// 4. 操作失败
	printResponse("操作失败", response.OperationFailed("发送邮件失败"))
	printResponse("操作失败(默认)", response.OperationFailed(""))

	// 5. 自定义业务码
	printResponse("自定义业务错误",
		response.Custom(10086, "积分不足，无法兑换", 200))
}

// ========== 实际场景测试 ==========

func TestRealWorldScenarios(t *testing.T) {
	fmt.Println("\n========== 实际业务场景 ==========")

	// 场景1: 用户注册
	fmt.Println("\n>>> 场景1: 用户注册")
	newUser := User{ID: 100, Username: "newuser", Email: "new@example.com"}

	// 成功
	printResponse("注册成功",
		response.Created(newUser).WithMessage("注册成功，欢迎加入"))

	// 失败 - 用户名已存在
	printResponse("注册失败 - 用户名已存在",
		response.RecordExists("用户名'newuser'已被注册"))

	// 场景2: 用户登录
	fmt.Println("\n>>> 场景2: 用户登录")
	loginData := map[string]any{
		"token":      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user":       User{ID: 1, Username: "张三", Email: "zhangsan@example.com"},
		"expires_in": 7200,
	}

	// 成功
	printResponse("登录成功",
		response.Data(loginData).WithMessage("登录成功"))

	// 失败 - 密码错误
	printResponse("登录失败 - 密码错误",
		response.Unauthorized("用户名或密码错误"))

	// 场景3: 获取用户详情
	fmt.Println("\n>>> 场景3: 获取用户详情")
	user := User{ID: 1, Username: "张三", Email: "zhangsan@example.com", Age: 25}

	// 成功
	printResponse("获取成功", response.Data(user))

	// 失败 - 用户不存在
	printResponse("获取失败 - 用户不存在",
		response.NotFound("用户ID不存在"))

	// 失败 - 未授权
	printResponse("获取失败 - 未授权",
		response.Unauthorized("请先登录"))

	// 场景4: 分页查询文章列表
	fmt.Println("\n>>> 场景4: 分页查询文章列表")
	articles := []Article{
		{ID: 1, Title: "Go语言入门", Content: "...", Author: "张三"},
		{ID: 2, Title: "微服务架构", Content: "...", Author: "李四"},
		{ID: 3, Title: "Docker实践", Content: "...", Author: "王五"},
	}

	printResponse("查询成功",
		response.Page(articles, 156, 1, 10).WithRequestID("req-001"))

	// 场景5: 更新用户信息
	fmt.Println("\n>>> 场景5: 更新用户信息")
	updatedUser := User{ID: 1, Username: "张三", Email: "new_email@example.com", Age: 26}

	// 成功
	printResponse("更新成功",
		response.Data(updatedUser).WithMessage("用户信息更新成功"))

	// 失败 - 权限不足
	printResponse("更新失败 - 权限不足",
		response.Forbidden("您只能修改自己的信息"))

	// 场景6: 删除文章
	fmt.Println("\n>>> 场景6: 删除文章")

	// 成功
	printResponse("删除成功", response.NoContent())

	// 失败 - 文章不存在
	printResponse("删除失败 - 文章不存在",
		response.NotFound("文章不存在或已被删除"))

	// 场景7: 支付订单
	fmt.Println("\n>>> 场景7: 支付订单")
	paymentResult := map[string]any{
		"order_id": "ORD20240101001",
		"amount":   99.99,
		"status":   "paid",
		"paid_at":  "2024-01-01T10:30:00Z",
	}

	// 成功
	printResponse("支付成功",
		response.Data(paymentResult).WithMessage("支付成功"))

	// 失败 - 余额不足
	printResponse("支付失败 - 余额不足",
		response.BusinessError("账户余额不足，请先充值"))

	// 失败 - 订单已支付
	printResponse("支付失败 - 订单已支付",
		response.BusinessError("该订单已支付，请勿重复操作"))

	// 场景8: 文件上传
	fmt.Println("\n>>> 场景8: 文件上传")
	uploadResult := map[string]any{
		"filename": "avatar.jpg",
		"url":      "https://cdn.example.com/uploads/avatar.jpg",
		"size":     1024000,
	}

	// 成功
	printResponse("上传成功",
		response.Created(uploadResult).WithMessage("文件上传成功"))

	// 失败 - 文件太大
	printResponse("上传失败 - 文件太大",
		response.BadRequest("文件大小超过限制(最大10MB)"))

	// 场景9: 批量操作
	fmt.Println("\n>>> 场景9: 批量删除用户")
	batchResult := map[string]any{
		"total":      10,
		"success":    8,
		"failed":     2,
		"failed_ids": []int{5, 9},
	}

	printResponse("批量操作结果",
		response.Data(batchResult).WithMessage("批量删除完成"))

	// 场景10: 系统错误
	fmt.Println("\n>>> 场景10: 系统错误")

	printResponse("数据库错误",
		response.InternalError("数据库连接失败，请稍后重试"))

	printResponse("服务降级",
		response.ServiceUnavailable("系统维护中，预计10分钟后恢复"))
}

// ========== 链式调用测试 ==========

func TestChainCalls(t *testing.T) {
	fmt.Println("\n========== 链式调用示例 ==========")

	user := User{ID: 1, Username: "张三", Email: "zhangsan@example.com"}

	// 完整链式调用
	printResponse("完整链式调用",
		response.Data(user).
			WithMessage("用户查询成功").
			WithRequestID("req-chain-001"))

	// 动态链式调用
	resp := response.Data(user)
	if true { // 模拟条件
		resp.WithMessage("VIP用户查询成功")
	}
	resp.WithRequestID("req-dynamic-001")
	printResponse("动态链式调用", resp)
}

// ========== 边界情况测试 ==========

func TestEdgeCases(t *testing.T) {
	fmt.Println("\n========== 边界情况测试 ==========")

	// nil 数据
	printResponse("nil数据", response.Data(nil))

	// 空切片
	printResponse("空切片", response.Data([]User{}))

	// 空分页
	printResponse("空分页", response.Page([]User{}, 0, 1, 10))

	// 空消息
	printResponse("空消息", response.BadRequest(""))

	// 零值结构体
	printResponse("零值结构体", response.Data(User{}))
}

// ========== 性能测试 ==========

func BenchmarkDataResponse(b *testing.B) {
	user := User{ID: 1, Username: "张三", Email: "zhangsan@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = response.Data(user)
	}
}

func BenchmarkPageResponse(b *testing.B) {
	users := []User{
		{ID: 1, Username: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Username: "李四", Email: "lisi@example.com"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = response.Page(users, 100, 1, 10)
	}
}

func BenchmarkChainCalls(b *testing.B) {
	user := User{ID: 1, Username: "张三", Email: "zhangsan@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = response.Data(user).
			WithMessage("查询成功").
			WithRequestID("req-001")
	}
}

// ========== 运行所有测试 ==========

func TestAll(t *testing.T) {
	TestSuccessResponses(t)
	TestErrorResponses(t)
	TestBusinessErrors(t)
	TestRealWorldScenarios(t)
	TestChainCalls(t)
	TestEdgeCases(t)
}
