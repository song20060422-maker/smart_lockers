// @ticket B端-中间件与鉴权单元测试
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 无 token
	r := gin.New()
	r.GET("/x", Auth(), func(c *gin.Context) { pkg.OK(c, nil) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
	if w.Body.String() == "" { t.Fatal("empty body") }

	// 带合法 token
	tok, _ := GenerateToken(7, model.RolePlatformOp, 0)
	r2 := gin.New()
	r2.GET("/x", Auth(), func(c *gin.Context) {
		u := User(c)
		if u.ID != 7 || u.Role != model.RolePlatformOp || u.Name() == "" { t.Errorf("u=%+v", u) }
		pkg.OK(c, nil)
	})
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2.ServeHTTP(w2, req2)

	// 错误 token
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/x", nil)
	req3.Header.Set("Authorization", "Bearer bad.token.here")
	r.ServeHTTP(w3, req3)
	if w3.Body.String() == "" { t.Fatal("bad token body empty") }
}

func TestRequirePermForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", func(c *gin.Context) { c.Set("role", model.RoleNormalUser); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	body := w.Body.String()
	if body == "" { t.Fatal("forbidden body empty") }
	// 放行路径
	r2 := gin.New()
	r2.GET("/x", func(c *gin.Context) { c.Set("role", model.RolePlatformOp); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })
	w2 := httptest.NewRecorder()
	r2.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w2.Body.String() == "" { t.Fatal("allowed body empty") }
}

func TestCommonMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(), RequestID(), Logger(), CORS())
	r.GET("/ok", func(c *gin.Context) { pkg.OK(c, map[string]string{"rid": c.GetString("request_id")}) })
	r.GET("/panic", func(c *gin.Context) { panic("boom") })
	// CORS 预检
	w0 := httptest.NewRecorder()
	req0 := httptest.NewRequest(http.MethodOptions, "/ok", nil)
	r.ServeHTTP(w0, req0)
	if w0.Header().Get("Access-Control-Allow-Origin") != "*" { t.Fatal("CORS header missing") }
	// RequestID 透传
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req1.Header.Set("X-Request-ID", "rid-xyz")
	r.ServeHTTP(w1, req1)
	if w1.Body.String() == "" || !contains(w1.Body.String(), "rid-xyz") { t.Fatalf("request id not propagated: %s", w1.Body.String()) }
	// Recovery
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if w2.Body.String() == "" { t.Fatal("recovery body empty") }
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ { if s[i:i+len(sub)] == sub { return true } }
	return false
}
