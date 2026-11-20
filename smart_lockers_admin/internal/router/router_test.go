// @ticket B端-路由集成单元测试
package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义函数 TestSetupRoutes：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestSetupRoutes(t *testing.T) {
	// 定义并初始化局部变量 r，取值为 Setup(store.New())

	r := Setup(store.New())
	// health
	// 定义并初始化局部变量 w，取值为 httptest.NewRecorder()

	w := httptest.NewRecorder()
	// 执行 r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))，完成对应的业务处理

	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
	// 声明变量 resp pkg.Response，用于保存运行过程中的数据

	var resp pkg.Response
	// 条件判断：err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.Code != 0 { t.Fatalf("health body=%s", w.Body.String()) } 时执行下面逻辑

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.Code != 0 { t.Fatalf("health body=%s", w.Body.String()) }
	// 未授权访问
	// 定义并初始化局部变量 w2，取值为 httptest.NewRecorder()

	w2 := httptest.NewRecorder()
	// 执行 r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/v1/admin/sites", nil))，完成对应的业务处理

	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/v1/admin/sites", nil))
	// 声明变量 resp2 pkg.Response，用于保存运行过程中的数据

	var resp2 pkg.Response
	// 执行 json.Unmarshal(w2.Body.Bytes(), &resp2)，完成对应的业务处理

	json.Unmarshal(w2.Body.Bytes(), &resp2)
	// 条件判断：resp2.Code != pkg.ErrAuth { t.Fatalf("unauth code=%d", resp2.Code) } 时执行下面逻辑

	if resp2.Code != pkg.ErrAuth { t.Fatalf("unauth code=%d", resp2.Code) }
	// 登录拿 token，再访问各功能（含越权校验）
	// 定义并初始化局部变量 loginBody，取值为 `{"username":"operator","password":"example"}`

	loginBody := `{"username":"operator","password":"example"}`
	// 定义并初始化局部变量 w3，取值为 httptest.NewRecorder()

	w3 := httptest.NewRecorder()
	// 定义并初始化局部变量 req3，取值为 httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))

	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	// 执行 req3.Header.Set("Content-Type", "application/json")，完成对应的业务处理

	req3.Header.Set("Content-Type", "application/json")
	// 执行 r.ServeHTTP(w3, req3)，完成对应的业务处理

	r.ServeHTTP(w3, req3)
	// 声明变量 lr pkg.Response，用于保存运行过程中的数据

	var lr pkg.Response
	// 执行 json.Unmarshal(w3.Body.Bytes(), &lr)，完成对应的业务处理

	json.Unmarshal(w3.Body.Bytes(), &lr)
	// 条件判断：lr.Code != 0 { t.Fatalf("login=%d body=%s", lr.Code, w3.Body.String()) } 时执行下面逻辑

	if lr.Code != 0 { t.Fatalf("login=%d body=%s", lr.Code, w3.Body.String()) }
	// 定义并初始化局部变量 data，取值为 lr.Data.(map[string]interface{})

	data := lr.Data.(map[string]interface{})
	// 定义并初始化局部变量 tok，取值为 data["token"].(string)

	tok := data["token"].(string)

	// 定义并初始化局部变量 authed，取值为 func(method, path, body string) pkg.Response {

	authed := func(method, path, body string) pkg.Response {
		// 定义并初始化局部变量 ww，取值为 httptest.NewRecorder()

		ww := httptest.NewRecorder()
		// 定义并初始化局部变量 rr，取值为 httptest.NewRequest(method, path, strings.NewReader(body))

		rr := httptest.NewRequest(method, path, strings.NewReader(body))
		// 执行 rr.Header.Set("Content-Type", "application/json")，完成对应的业务处理

		rr.Header.Set("Content-Type", "application/json")
		// 执行 rr.Header.Set("Authorization", "Bearer "+tok)，完成对应的业务处理

		rr.Header.Set("Authorization", "Bearer "+tok)
		// 执行 r.ServeHTTP(ww, rr)，完成对应的业务处理

		r.ServeHTTP(ww, rr)
		// 声明变量 x pkg.Response，用于保存运行过程中的数据

		var x pkg.Response
		// 执行 json.Unmarshal(ww.Body.Bytes(), &x)，完成对应的业务处理

		json.Unmarshal(ww.Body.Bytes(), &x)
		// 返回 x 给调用方，结束当前函数

		return x
	// 当前代码块结束

	}
	// 循环处理：for _, p := range []string{"/api/v1/admin/sites", "/api/v1/admin/devices", "/api/v1/admin/devices/monitor", "/api/v1/admin/reconcile", "/api/v1/admin/alarms"...，逐轮迭代

	for _, p := range []string{"/api/v1/admin/sites", "/api/v1/admin/devices", "/api/v1/admin/devices/monitor", "/api/v1/admin/reconcile", "/api/v1/admin/alarms", "/api/v1/auth/me"} {
		// 条件判断：x := authed(http.MethodGet, p, ""); x.Code != 0 { t.Fatalf("GET %s code=%d", p, x.Code) } 时执行下面逻辑

		if x := authed(http.MethodGet, p, ""); x.Code != 0 { t.Fatalf("GET %s code=%d", p, x.Code) }
	// 当前代码块结束

	}
	// 条件判断：x := authed(http.MethodPost, "/api/v1/admin/reconcile/run", `{"batch_no":"B-rt","expected_amount":1,"actual_amount":1}`); x.Code != 0 { t.Fatalf("run code=%d... 时执行下面逻辑

	if x := authed(http.MethodPost, "/api/v1/admin/reconcile/run", `{"batch_no":"B-rt","expected_amount":1,"actual_amount":1}`); x.Code != 0 { t.Fatalf("run code=%d", x.Code) }
// 当前代码块结束

}
