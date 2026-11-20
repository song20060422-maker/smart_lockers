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

// 定义函数 TestAuthMiddleware：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestAuthMiddleware(t *testing.T) {
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 无 token
	// 定义并初始化局部变量 r，取值为 gin.New()

	r := gin.New()
	// 执行 r.GET("/x", Auth(), func(c *gin.Context) { pkg.OK(c, nil) })，完成对应的业务处理

	r.GET("/x", Auth(), func(c *gin.Context) { pkg.OK(c, nil) })
	// 定义并初始化局部变量 w，取值为 httptest.NewRecorder()

	w := httptest.NewRecorder()
	// 执行 r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))，完成对应的业务处理

	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	// 条件判断：w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) } 时执行下面逻辑

	if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
	// 条件判断：w.Body.String() == "" { t.Fatal("empty body") } 时执行下面逻辑

	if w.Body.String() == "" { t.Fatal("empty body") }

	// 带合法 token
	// 定义并初始化局部变量 tok, _，取值为 GenerateToken(7, model.RolePlatformOp, 0)

	tok, _ := GenerateToken(7, model.RolePlatformOp, 0)
	// 定义并初始化局部变量 r2，取值为 gin.New()

	r2 := gin.New()
	// 执行 r2.GET("/x", Auth(), func(c *gin.Context) {，完成对应的业务处理

	r2.GET("/x", Auth(), func(c *gin.Context) {
		// 定义并初始化局部变量 u，取值为 User(c)

		u := User(c)
		// 条件判断：u.ID != 7 || u.Role != model.RolePlatformOp || u.Name() == "" { t.Errorf("u=%+v", u) } 时执行下面逻辑

		if u.ID != 7 || u.Role != model.RolePlatformOp || u.Name() == "" { t.Errorf("u=%+v", u) }
		// 执行 pkg.OK(c, nil)，完成对应的业务处理

		pkg.OK(c, nil)
	// 执行 })，完成对应的业务处理

	})
	// 定义并初始化局部变量 w2，取值为 httptest.NewRecorder()

	w2 := httptest.NewRecorder()
	// 定义并初始化局部变量 req2，取值为 httptest.NewRequest(http.MethodGet, "/x", nil)

	req2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	// 执行 req2.Header.Set("Authorization", "Bearer "+tok)，完成对应的业务处理

	req2.Header.Set("Authorization", "Bearer "+tok)
	// 执行 r2.ServeHTTP(w2, req2)，完成对应的业务处理

	r2.ServeHTTP(w2, req2)

	// 错误 token
	// 定义并初始化局部变量 w3，取值为 httptest.NewRecorder()

	w3 := httptest.NewRecorder()
	// 定义并初始化局部变量 req3，取值为 httptest.NewRequest(http.MethodGet, "/x", nil)

	req3 := httptest.NewRequest(http.MethodGet, "/x", nil)
	// 执行 req3.Header.Set("Authorization", "Bearer bad.token.here")，完成对应的业务处理

	req3.Header.Set("Authorization", "Bearer bad.token.here")
	// 执行 r.ServeHTTP(w3, req3)，完成对应的业务处理

	r.ServeHTTP(w3, req3)
	// 条件判断：w3.Body.String() == "" { t.Fatal("bad token body empty") } 时执行下面逻辑

	if w3.Body.String() == "" { t.Fatal("bad token body empty") }
// 当前代码块结束

}

// 定义函数 TestRequirePermForbidden：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestRequirePermForbidden(t *testing.T) {
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 定义并初始化局部变量 r，取值为 gin.New()

	r := gin.New()
	// 执行 r.GET("/x", func(c *gin.Context) { c.Set("role", model.RoleNormalUser); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })，完成对应的业务处理

	r.GET("/x", func(c *gin.Context) { c.Set("role", model.RoleNormalUser); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })
	// 定义并初始化局部变量 w，取值为 httptest.NewRecorder()

	w := httptest.NewRecorder()
	// 执行 r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))，完成对应的业务处理

	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	// 定义并初始化局部变量 body，取值为 w.Body.String()

	body := w.Body.String()
	// 条件判断：body == "" { t.Fatal("forbidden body empty") } 时执行下面逻辑

	if body == "" { t.Fatal("forbidden body empty") }
	// 放行路径
	// 定义并初始化局部变量 r2，取值为 gin.New()

	r2 := gin.New()
	// 执行 r2.GET("/x", func(c *gin.Context) { c.Set("role", model.RolePlatformOp); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })，完成对应的业务处理

	r2.GET("/x", func(c *gin.Context) { c.Set("role", model.RolePlatformOp); c.Next() }, RequirePerm(PermReconRun), func(c *gin.Context) { pkg.OK(c, nil) })
	// 定义并初始化局部变量 w2，取值为 httptest.NewRecorder()

	w2 := httptest.NewRecorder()
	// 执行 r2.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/x", nil))，完成对应的业务处理

	r2.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/x", nil))
	// 条件判断：w2.Body.String() == "" { t.Fatal("allowed body empty") } 时执行下面逻辑

	if w2.Body.String() == "" { t.Fatal("allowed body empty") }
// 当前代码块结束

}

// 定义函数 TestCommonMiddlewares：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCommonMiddlewares(t *testing.T) {
	// 执行 gin.SetMode(gin.TestMode)，完成对应的业务处理

	gin.SetMode(gin.TestMode)
	// 定义并初始化局部变量 r，取值为 gin.New()

	r := gin.New()
	// 执行 r.Use(Recovery(), RequestID(), Logger(), CORS())，完成对应的业务处理

	r.Use(Recovery(), RequestID(), Logger(), CORS())
	// 执行 r.GET("/ok", func(c *gin.Context) { pkg.OK(c, map[string]string{"rid": c.GetString("request_id")}) })，完成对应的业务处理

	r.GET("/ok", func(c *gin.Context) { pkg.OK(c, map[string]string{"rid": c.GetString("request_id")}) })
	// 执行 r.GET("/panic", func(c *gin.Context) { panic("boom") })，完成对应的业务处理

	r.GET("/panic", func(c *gin.Context) { panic("boom") })
	// CORS 预检
	// 定义并初始化局部变量 w0，取值为 httptest.NewRecorder()

	w0 := httptest.NewRecorder()
	// 定义并初始化局部变量 req0，取值为 httptest.NewRequest(http.MethodOptions, "/ok", nil)

	req0 := httptest.NewRequest(http.MethodOptions, "/ok", nil)
	// 执行 r.ServeHTTP(w0, req0)，完成对应的业务处理

	r.ServeHTTP(w0, req0)
	// 条件判断：w0.Header().Get("Access-Control-Allow-Origin") != "*" { t.Fatal("CORS header missing") } 时执行下面逻辑

	if w0.Header().Get("Access-Control-Allow-Origin") != "*" { t.Fatal("CORS header missing") }
	// RequestID 透传
	// 定义并初始化局部变量 w1，取值为 httptest.NewRecorder()

	w1 := httptest.NewRecorder()
	// 定义并初始化局部变量 req1，取值为 httptest.NewRequest(http.MethodGet, "/ok", nil)

	req1 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	// 执行 req1.Header.Set("X-Request-ID", "rid-xyz")，完成对应的业务处理

	req1.Header.Set("X-Request-ID", "rid-xyz")
	// 执行 r.ServeHTTP(w1, req1)，完成对应的业务处理

	r.ServeHTTP(w1, req1)
	// 条件判断：w1.Body.String() == "" || !contains(w1.Body.String(), "rid-xyz") { t.Fatalf("request id not propagated: %s", w1.Body.String()) } 时执行下面逻辑

	if w1.Body.String() == "" || !contains(w1.Body.String(), "rid-xyz") { t.Fatalf("request id not propagated: %s", w1.Body.String()) }
	// Recovery
	// 定义并初始化局部变量 w2，取值为 httptest.NewRecorder()

	w2 := httptest.NewRecorder()
	// 执行 r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/panic", nil))，完成对应的业务处理

	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/panic", nil))
	// 条件判断：w2.Body.String() == "" { t.Fatal("recovery body empty") } 时执行下面逻辑

	if w2.Body.String() == "" { t.Fatal("recovery body empty") }
// 当前代码块结束

}

// 定义函数 contains：封装一段可复用的处理逻辑（参数：s, sub string）

func contains(s, sub string) bool {
	// 循环处理：for i := 0; i+len(sub) <= len(s); i++ { if s[i:i+len(sub)] == sub { return true } }，逐轮迭代

	for i := 0; i+len(sub) <= len(s); i++ { if s[i:i+len(sub)] == sub { return true } }
	// 返回 false 给调用方，结束当前函数

	return false
// 当前代码块结束

}
