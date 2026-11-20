// @ticket B端-HTTP handler 单元测试（覆盖四大功能入口与参数校验）
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/store"
)

// 定义函数 ctx：封装一段可复用的处理逻辑（参数：method, path, body string, role model.AdminRole, siteID uint）

func ctx(method, path, body string, role model.AdminRole, siteID uint) (*gin.Context, *httptest.ResponseRecorder) {
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 定义并初始化局部变量 w，取值为 httptest.NewRecorder()

	w := httptest.NewRecorder()
	// 定义并初始化局部变量 c, _，取值为 gin.CreateTestContext(w)

	c, _ := gin.CreateTestContext(w)
	// 为 c.Request 赋值：httptest.NewRequest(method, path, strings.NewReader(body))

	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	// 执行 c.Request.Header.Set("Content-Type", "application/json")，完成对应的业务处理

	c.Request.Header.Set("Content-Type", "application/json")
	// 执行 c.Set("uid", uint(1)); c.Set("role", role); c.Set("site_id", siteID)，完成对应的业务处理

	c.Set("uid", uint(1)); c.Set("role", role); c.Set("site_id", siteID)
	// 返回 c, w 给调用方，结束当前函数

	return c, w
// 当前代码块结束

}

// 定义函数 decode：封装一段可复用的处理逻辑（参数：t *testing.T, w *httptest.ResponseRecorder）

func decode(t *testing.T, w *httptest.ResponseRecorder) pkg.Response {
	// 声明变量 r pkg.Response，用于保存运行过程中的数据

	var r pkg.Response
	// 条件判断：err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatalf("decode %v body=%s", err, w.Body.String()) } 时执行下面逻辑

	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatalf("decode %v body=%s", err, w.Body.String()) }
	// 返回 r 给调用方，结束当前函数

	return r
// 当前代码块结束

}

// 定义函数 TestAuthHandlers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestAuthHandlers(t *testing.T) {
	// 定义并初始化局部变量 h，取值为 &AuthHandler{}

	h := &AuthHandler{}
	// 登录成功（默认占位账号）
	// 定义并初始化局部变量 c, w，取值为 ctx(http.MethodPost, "/login", `{"username":"operator","password":"example"}`, model.RoleNormalUser, 0)

	c, w := ctx(http.MethodPost, "/login", `{"username":"operator","password":"example"}`, model.RoleNormalUser, 0)
	// 执行 h.Login(c)，完成对应的业务处理

	h.Login(c)
	// 条件判断：r := decode(t, w); r.Code != 0 { t.Fatalf("login code=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w); r.Code != 0 { t.Fatalf("login code=%d", r.Code) }
	// 登录失败
	// 定义并初始化局部变量 c2, w2，取值为 ctx(http.MethodPost, "/login", `{"username":"x","password":"y"}`, model.RoleNormalUser, 0)

	c2, w2 := ctx(http.MethodPost, "/login", `{"username":"x","password":"y"}`, model.RoleNormalUser, 0)
	// 执行 h.Login(c2)，完成对应的业务处理

	h.Login(c2)
	// 条件判断：r := decode(t, w2); r.Code == 0 { t.Fatal("bad login should fail") } 时执行下面逻辑

	if r := decode(t, w2); r.Code == 0 { t.Fatal("bad login should fail") }
	// 绑定错误
	// 定义并初始化局部变量 c3, w3，取值为 ctx(http.MethodPost, "/login", `{`, model.RoleNormalUser, 0)

	c3, w3 := ctx(http.MethodPost, "/login", `{`, model.RoleNormalUser, 0)
	// 执行 h.Login(c3)，完成对应的业务处理

	h.Login(c3)
	// 条件判断：r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bind code=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bind code=%d", r.Code) }
	// Me
	// 定义并初始化局部变量 c4, w4，取值为 ctx(http.MethodGet, "/me", "", model.RoleSiteAdmin, 1001)

	c4, w4 := ctx(http.MethodGet, "/me", "", model.RoleSiteAdmin, 1001)
	// 执行 h.Me(c4)，完成对应的业务处理

	h.Me(c4)
	// 条件判断：r := decode(t, w4); r.Code != 0 { t.Fatalf("me code=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w4); r.Code != 0 { t.Fatalf("me code=%d", r.Code) }
// 当前代码块结束

}

// 定义函数 TestSiteHandlers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestSiteHandlers(t *testing.T) {
	// 定义并初始化局部变量 h，取值为 NewSite(service.NewSite(store.New()))

	h := NewSite(service.NewSite(store.New()))
	// 定义并初始化局部变量 op，取值为 model.RolePlatformOp

	op := model.RolePlatformOp
	// 定义并初始化局部变量 c, w，取值为 ctx(http.MethodGet, "/sites", "", op, 0); h.List(c)

	c, w := ctx(http.MethodGet, "/sites", "", op, 0); h.List(c)
	// 条件判断：r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	// 定义并初始化局部变量 c2, w2，取值为 ctx(http.MethodPost, "/sites", `{"code":"S9","name":"n"}`, op, 0); h.Create(c2)

	c2, w2 := ctx(http.MethodPost, "/sites", `{"code":"S9","name":"n"}`, op, 0); h.Create(c2)
	// 条件判断：r := decode(t, w2); r.Code != 0 { t.Fatalf("create=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w2); r.Code != 0 { t.Fatalf("create=%d", r.Code) }
	// 定义并初始化局部变量 c3, w3，取值为 ctx(http.MethodPost, "/sites", `{`, op, 0); h.Create(c3)

	c3, w3 := ctx(http.MethodPost, "/sites", `{`, op, 0); h.Create(c3)
	// 条件判断：r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad create=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad create=%d", r.Code) }
	// 定义并初始化局部变量 c4, w4，取值为 ctx(http.MethodPost, "/sites/1001/status", `{"status":0}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "1001"}}; h.UpdateStatus(c4)

	c4, w4 := ctx(http.MethodPost, "/sites/1001/status", `{"status":0}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "1001"}}; h.UpdateStatus(c4)
	// 条件判断：r := decode(t, w4); r.Code != 0 { t.Fatalf("status=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w4); r.Code != 0 { t.Fatalf("status=%d", r.Code) }
	// 定义并初始化局部变量 c5, w5，取值为 ctx(http.MethodPost, "/sites/abc/status", `{"status":0}`, op, 0); c5.Params = gin.Params{{Key: "id", Value: "abc"}}; h.UpdateStatus(c5)

	c5, w5 := ctx(http.MethodPost, "/sites/abc/status", `{"status":0}`, op, 0); c5.Params = gin.Params{{Key: "id", Value: "abc"}}; h.UpdateStatus(c5)
	// 条件判断：r := decode(t, w5); r.Code != pkg.ErrBad { t.Fatalf("bad id=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w5); r.Code != pkg.ErrBad { t.Fatalf("bad id=%d", r.Code) }
// 当前代码块结束

}

// 定义函数 TestDeviceHandlers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestDeviceHandlers(t *testing.T) {
	// 定义并初始化局部变量 h，取值为 NewDevice(service.NewDevice(store.New()))

	h := NewDevice(service.NewDevice(store.New()))
	// 定义并初始化局部变量 op，取值为 model.RolePlatformOp

	op := model.RolePlatformOp
	// 定义并初始化局部变量 c, w，取值为 ctx(http.MethodGet, "/devices", "", op, 0); h.List(c)

	c, w := ctx(http.MethodGet, "/devices", "", op, 0); h.List(c)
	// 条件判断：r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	// 定义并初始化局部变量 c2, w2，取值为 ctx(http.MethodGet, "/devices/monitor", "", op, 0); h.Monitor(c2)

	c2, w2 := ctx(http.MethodGet, "/devices/monitor", "", op, 0); h.Monitor(c2)
	// 条件判断：r := decode(t, w2); r.Code != 0 { t.Fatalf("monitor=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w2); r.Code != 0 { t.Fatalf("monitor=%d", r.Code) }
	// 定义并初始化局部变量 c3, w3，取值为 ctx(http.MethodPost, "/devices/2001/control", `{"status":1}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "2001"}}; h.Control(c3)

	c3, w3 := ctx(http.MethodPost, "/devices/2001/control", `{"status":1}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "2001"}}; h.Control(c3)
	// 条件判断：r := decode(t, w3); r.Code != 0 { t.Fatalf("control=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w3); r.Code != 0 { t.Fatalf("control=%d", r.Code) }
	// 定义并初始化局部变量 c4, w4，取值为 ctx(http.MethodPost, "/devices/abc/control", `{"status":1}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Control(c4)

	c4, w4 := ctx(http.MethodPost, "/devices/abc/control", `{"status":1}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Control(c4)
	// 条件判断：r := decode(t, w4); r.Code != pkg.ErrBad { t.Fatalf("bad control id=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w4); r.Code != pkg.ErrBad { t.Fatalf("bad control id=%d", r.Code) }
// 当前代码块结束

}

// 定义函数 TestReconcileHandlers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestReconcileHandlers(t *testing.T) {
	// 定义并初始化局部变量 h，取值为 NewReconcile(service.NewReconcile(store.New()))

	h := NewReconcile(service.NewReconcile(store.New()))
	// 定义并初始化局部变量 op，取值为 model.RolePlatformOp

	op := model.RolePlatformOp
	// 定义并初始化局部变量 c, w，取值为 ctx(http.MethodGet, "/reconcile", "", op, 0); h.List(c)

	c, w := ctx(http.MethodGet, "/reconcile", "", op, 0); h.List(c)
	// 条件判断：r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	// 定义并初始化局部变量 c2, w2，取值为 ctx(http.MethodPost, "/reconcile/run", `{"batch_no":"B-h","expected_amount":100,"actual_amount":100}`, op, 0); h.Run(c2)

	c2, w2 := ctx(http.MethodPost, "/reconcile/run", `{"batch_no":"B-h","expected_amount":100,"actual_amount":100}`, op, 0); h.Run(c2)
	// 条件判断：r := decode(t, w2); r.Code != 0 { t.Fatalf("run=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w2); r.Code != 0 { t.Fatalf("run=%d", r.Code) }
	// 定义并初始化局部变量 c3, w3，取值为 ctx(http.MethodPost, "/reconcile/run", `{`, op, 0); h.Run(c3)

	c3, w3 := ctx(http.MethodPost, "/reconcile/run", `{`, op, 0); h.Run(c3)
	// 条件判断：r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad run=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad run=%d", r.Code) }
// 当前代码块结束

}

// 定义函数 TestAlarmHandlers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestAlarmHandlers(t *testing.T) {
	// 定义并初始化局部变量 h，取值为 NewAlarm(service.NewAlarm(store.New()))

	h := NewAlarm(service.NewAlarm(store.New()))
	// 定义并初始化局部变量 op，取值为 model.RolePlatformOp

	op := model.RolePlatformOp
	// 定义并初始化局部变量 c, w，取值为 ctx(http.MethodGet, "/alarms", "", op, 0); h.List(c)

	c, w := ctx(http.MethodGet, "/alarms", "", op, 0); h.List(c)
	// 条件判断：r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	// 定义并初始化局部变量 c2, w2，取值为 ctx(http.MethodPost, "/alarms/4001/handle", `{"note":"ok"}`, op, 0); c2.Params = gin.Params{{Key: "id", Value: "4001"}}; h.Handle(c2)

	c2, w2 := ctx(http.MethodPost, "/alarms/4001/handle", `{"note":"ok"}`, op, 0); c2.Params = gin.Params{{Key: "id", Value: "4001"}}; h.Handle(c2)
	// 条件判断：r := decode(t, w2); r.Code != 0 { t.Fatalf("handle=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w2); r.Code != 0 { t.Fatalf("handle=%d", r.Code) }
	// 定义并初始化局部变量 c3, w3，取值为 ctx(http.MethodPost, "/alarms/abc/handle", `{}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Handle(c3)

	c3, w3 := ctx(http.MethodPost, "/alarms/abc/handle", `{}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Handle(c3)
	// 条件判断：r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad handle=%d", r.Code) } 时执行下面逻辑

	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad handle=%d", r.Code) }
// 当前代码块结束

}

// 定义函数 TestHelpers：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestHelpers(t *testing.T) {
	// 条件判断：!bytes.Equal([]byte("x"), []byte("x")) { t.Fatal() } 时执行下面逻辑

	if !bytes.Equal([]byte("x"), []byte("x")) { t.Fatal() }
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 定义并初始化局部变量 c, _，取值为 gin.CreateTestContext(httptest.NewRecorder())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	// 为 c.Request 赋值：httptest.NewRequest(http.MethodGet, "/?idempotency_key=k", nil)

	c.Request = httptest.NewRequest(http.MethodGet, "/?idempotency_key=k", nil)
	// 条件判断：idemKey(c) != "k" { t.Fatal("query idem key") } 时执行下面逻辑

	if idemKey(c) != "k" { t.Fatal("query idem key") }
	// 执行 c.Request.Header.Set("Idempotency-Key", "hdr")，完成对应的业务处理

	c.Request.Header.Set("Idempotency-Key", "hdr")
	// 条件判断：idemKey(c) != "hdr" { t.Fatal("header idem key priority") } 时执行下面逻辑

	if idemKey(c) != "hdr" { t.Fatal("header idem key priority") }
	// 为 c.Params 赋值：gin.Params{{Key: "id", Value: "42"}}

	c.Params = gin.Params{{Key: "id", Value: "42"}}
	// 条件判断：v, ok := uintParam(c, "id"); !ok || v != 42 { t.Fatal("uintParam ok") } 时执行下面逻辑

	if v, ok := uintParam(c, "id"); !ok || v != 42 { t.Fatal("uintParam ok") }
	// 为 c.Params 赋值：gin.Params{{Key: "id", Value: "0"}}

	c.Params = gin.Params{{Key: "id", Value: "0"}}
	// 条件判断：_, ok := uintParam(c, "id"); ok { t.Fatal("zero id should fail") } 时执行下面逻辑

	if _, ok := uintParam(c, "id"); ok { t.Fatal("zero id should fail") }
	// 条件判断：v := middleware.SiteScope(middleware.CurrentUser{Role: model.RoleNormalUser}); v != 0 { t.Fatalf("scope=%d", v) } 时执行下面逻辑

	if v := middleware.SiteScope(middleware.CurrentUser{Role: model.RoleNormalUser}); v != 0 { t.Fatalf("scope=%d", v) }
// 当前代码块结束

}
