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

func ctx(method, path, body string, role model.AdminRole, siteID uint) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("uid", uint(1)); c.Set("role", role); c.Set("site_id", siteID)
	return c, w
}

func decode(t *testing.T, w *httptest.ResponseRecorder) pkg.Response {
	var r pkg.Response
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatalf("decode %v body=%s", err, w.Body.String()) }
	return r
}

func TestAuthHandlers(t *testing.T) {
	h := &AuthHandler{}
	// 登录成功（默认占位账号）
	c, w := ctx(http.MethodPost, "/login", `{"username":"operator","password":"example"}`, model.RoleNormalUser, 0)
	h.Login(c)
	if r := decode(t, w); r.Code != 0 { t.Fatalf("login code=%d", r.Code) }
	// 登录失败
	c2, w2 := ctx(http.MethodPost, "/login", `{"username":"x","password":"y"}`, model.RoleNormalUser, 0)
	h.Login(c2)
	if r := decode(t, w2); r.Code == 0 { t.Fatal("bad login should fail") }
	// 绑定错误
	c3, w3 := ctx(http.MethodPost, "/login", `{`, model.RoleNormalUser, 0)
	h.Login(c3)
	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bind code=%d", r.Code) }
	// Me
	c4, w4 := ctx(http.MethodGet, "/me", "", model.RoleSiteAdmin, 1001)
	h.Me(c4)
	if r := decode(t, w4); r.Code != 0 { t.Fatalf("me code=%d", r.Code) }
}

func TestSiteHandlers(t *testing.T) {
	h := NewSite(service.NewSite(store.New()))
	op := model.RolePlatformOp
	c, w := ctx(http.MethodGet, "/sites", "", op, 0); h.List(c)
	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	c2, w2 := ctx(http.MethodPost, "/sites", `{"code":"S9","name":"n"}`, op, 0); h.Create(c2)
	if r := decode(t, w2); r.Code != 0 { t.Fatalf("create=%d", r.Code) }
	c3, w3 := ctx(http.MethodPost, "/sites", `{`, op, 0); h.Create(c3)
	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad create=%d", r.Code) }
	c4, w4 := ctx(http.MethodPost, "/sites/1001/status", `{"status":0}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "1001"}}; h.UpdateStatus(c4)
	if r := decode(t, w4); r.Code != 0 { t.Fatalf("status=%d", r.Code) }
	c5, w5 := ctx(http.MethodPost, "/sites/abc/status", `{"status":0}`, op, 0); c5.Params = gin.Params{{Key: "id", Value: "abc"}}; h.UpdateStatus(c5)
	if r := decode(t, w5); r.Code != pkg.ErrBad { t.Fatalf("bad id=%d", r.Code) }
}

func TestDeviceHandlers(t *testing.T) {
	h := NewDevice(service.NewDevice(store.New()))
	op := model.RolePlatformOp
	c, w := ctx(http.MethodGet, "/devices", "", op, 0); h.List(c)
	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	c2, w2 := ctx(http.MethodGet, "/devices/monitor", "", op, 0); h.Monitor(c2)
	if r := decode(t, w2); r.Code != 0 { t.Fatalf("monitor=%d", r.Code) }
	c3, w3 := ctx(http.MethodPost, "/devices/2001/control", `{"status":1}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "2001"}}; h.Control(c3)
	if r := decode(t, w3); r.Code != 0 { t.Fatalf("control=%d", r.Code) }
	c4, w4 := ctx(http.MethodPost, "/devices/abc/control", `{"status":1}`, op, 0); c4.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Control(c4)
	if r := decode(t, w4); r.Code != pkg.ErrBad { t.Fatalf("bad control id=%d", r.Code) }
}

func TestReconcileHandlers(t *testing.T) {
	h := NewReconcile(service.NewReconcile(store.New()))
	op := model.RolePlatformOp
	c, w := ctx(http.MethodGet, "/reconcile", "", op, 0); h.List(c)
	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	c2, w2 := ctx(http.MethodPost, "/reconcile/run", `{"batch_no":"B-h","expected_amount":100,"actual_amount":100}`, op, 0); h.Run(c2)
	if r := decode(t, w2); r.Code != 0 { t.Fatalf("run=%d", r.Code) }
	c3, w3 := ctx(http.MethodPost, "/reconcile/run", `{`, op, 0); h.Run(c3)
	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad run=%d", r.Code) }
}

func TestAlarmHandlers(t *testing.T) {
	h := NewAlarm(service.NewAlarm(store.New()))
	op := model.RolePlatformOp
	c, w := ctx(http.MethodGet, "/alarms", "", op, 0); h.List(c)
	if r := decode(t, w); r.Code != 0 { t.Fatalf("list=%d", r.Code) }
	c2, w2 := ctx(http.MethodPost, "/alarms/4001/handle", `{"note":"ok"}`, op, 0); c2.Params = gin.Params{{Key: "id", Value: "4001"}}; h.Handle(c2)
	if r := decode(t, w2); r.Code != 0 { t.Fatalf("handle=%d", r.Code) }
	c3, w3 := ctx(http.MethodPost, "/alarms/abc/handle", `{}`, op, 0); c3.Params = gin.Params{{Key: "id", Value: "abc"}}; h.Handle(c3)
	if r := decode(t, w3); r.Code != pkg.ErrBad { t.Fatalf("bad handle=%d", r.Code) }
}

func TestHelpers(t *testing.T) {
	if !bytes.Equal([]byte("x"), []byte("x")) { t.Fatal() }
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?idempotency_key=k", nil)
	if idemKey(c) != "k" { t.Fatal("query idem key") }
	c.Request.Header.Set("Idempotency-Key", "hdr")
	if idemKey(c) != "hdr" { t.Fatal("header idem key priority") }
	c.Params = gin.Params{{Key: "id", Value: "42"}}
	if v, ok := uintParam(c, "id"); !ok || v != 42 { t.Fatal("uintParam ok") }
	c.Params = gin.Params{{Key: "id", Value: "0"}}
	if _, ok := uintParam(c, "id"); ok { t.Fatal("zero id should fail") }
	if v := middleware.SiteScope(middleware.CurrentUser{Role: model.RoleNormalUser}); v != 0 { t.Fatalf("scope=%d", v) }
}
