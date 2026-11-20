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

func TestSetupRoutes(t *testing.T) {
	r := Setup(store.New())
	// health
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
	var resp pkg.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.Code != 0 { t.Fatalf("health body=%s", w.Body.String()) }
	// 未授权访问
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/v1/admin/sites", nil))
	var resp2 pkg.Response
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	if resp2.Code != pkg.ErrAuth { t.Fatalf("unauth code=%d", resp2.Code) }
	// 登录拿 token，再访问各功能（含越权校验）
	loginBody := `{"username":"operator","password":"example"}`
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	req3.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w3, req3)
	var lr pkg.Response
	json.Unmarshal(w3.Body.Bytes(), &lr)
	if lr.Code != 0 { t.Fatalf("login=%d body=%s", lr.Code, w3.Body.String()) }
	data := lr.Data.(map[string]interface{})
	tok := data["token"].(string)

	authed := func(method, path, body string) pkg.Response {
		ww := httptest.NewRecorder()
		rr := httptest.NewRequest(method, path, strings.NewReader(body))
		rr.Header.Set("Content-Type", "application/json")
		rr.Header.Set("Authorization", "Bearer "+tok)
		r.ServeHTTP(ww, rr)
		var x pkg.Response
		json.Unmarshal(ww.Body.Bytes(), &x)
		return x
	}
	for _, p := range []string{"/api/v1/admin/sites", "/api/v1/admin/devices", "/api/v1/admin/devices/monitor", "/api/v1/admin/reconcile", "/api/v1/admin/alarms", "/api/v1/auth/me"} {
		if x := authed(http.MethodGet, p, ""); x.Code != 0 { t.Fatalf("GET %s code=%d", p, x.Code) }
	}
	if x := authed(http.MethodPost, "/api/v1/admin/reconcile/run", `{"batch_no":"B-rt","expected_amount":1,"actual_amount":1}`); x.Code != 0 { t.Fatalf("run code=%d", x.Code) }
}
