// @ticket B端-统一响应单元测试
package pkg

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 定义函数 newCtx：封装一段可复用的处理逻辑（参数：）

func newCtx() (*gin.Context, *httptest.ResponseRecorder) {
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 定义并初始化局部变量 w，取值为 httptest.NewRecorder()

	w := httptest.NewRecorder()
	// 定义并初始化局部变量 c, _，取值为 gin.CreateTestContext(w)

	c, _ := gin.CreateTestContext(w)
	// 返回 c, w 给调用方，结束当前函数

	return c, w
// 当前代码块结束

}

// 定义函数 TestResponseOKAndFail：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestResponseOKAndFail(t *testing.T) {
	// 定义并初始化局部变量 c, w，取值为 newCtx()

	c, w := newCtx()
	// 执行 c.Set("request_id", "rid-1")，完成对应的业务处理

	c.Set("request_id", "rid-1")
	// 执行 OK(c, map[string]int{"n": 1})，完成对应的业务处理

	OK(c, map[string]int{"n": 1})
	// 声明变量 r Response，用于保存运行过程中的数据

	var r Response
	// 条件判断：err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatal(err) } 时执行下面逻辑

	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatal(err) }
	// 条件判断：r.Code != 0 || r.Message != "ok" || r.RequestID != "rid-1" { t.Fatalf("OK resp=%+v", r) } 时执行下面逻辑

	if r.Code != 0 || r.Message != "ok" || r.RequestID != "rid-1" { t.Fatalf("OK resp=%+v", r) }

	// 定义并初始化局部变量 c2, w2，取值为 newCtx()

	c2, w2 := newCtx()
	// 执行 Fail(c2, ErrForb)，完成对应的业务处理

	Fail(c2, ErrForb)
	// 声明变量 r2 Response，用于保存运行过程中的数据

	var r2 Response
	// 执行 json.Unmarshal(w2.Body.Bytes(), &r2)，完成对应的业务处理

	json.Unmarshal(w2.Body.Bytes(), &r2)
	// 条件判断：r2.Code != ErrForb || r2.Message == "" { t.Fatalf("Fail resp=%+v", r2) } 时执行下面逻辑

	if r2.Code != ErrForb || r2.Message == "" { t.Fatalf("Fail resp=%+v", r2) }

	// 定义并初始化局部变量 c3, w3，取值为 newCtx()

	c3, w3 := newCtx()
	// 执行 FailMsg(c3, ErrBad, "custom")，完成对应的业务处理

	FailMsg(c3, ErrBad, "custom")
	// 声明变量 r3 Response，用于保存运行过程中的数据

	var r3 Response
	// 执行 json.Unmarshal(w3.Body.Bytes(), &r3)，完成对应的业务处理

	json.Unmarshal(w3.Body.Bytes(), &r3)
	// 条件判断：r3.Message != "custom" { t.Fatalf("FailMsg resp=%+v", r3) } 时执行下面逻辑

	if r3.Message != "custom" { t.Fatalf("FailMsg resp=%+v", r3) }

	// 定义并初始化局部变量 c4, _，取值为 newCtx()

	c4, _ := newCtx()
	// 执行 AbortWithCode(c4, ErrAuth)，完成对应的业务处理

	AbortWithCode(c4, ErrAuth)
	// 条件判断：!c4.IsAborted() { t.Fatal("AbortWithCode did not abort") } 时执行下面逻辑

	if !c4.IsAborted() { t.Fatal("AbortWithCode did not abort") }
// 当前代码块结束

}
