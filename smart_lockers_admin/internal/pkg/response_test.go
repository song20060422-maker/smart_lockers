// @ticket B端-统一响应单元测试
package pkg

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newCtx() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestResponseOKAndFail(t *testing.T) {
	c, w := newCtx()
	c.Set("request_id", "rid-1")
	OK(c, map[string]int{"n": 1})
	var r Response
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil { t.Fatal(err) }
	if r.Code != 0 || r.Message != "ok" || r.RequestID != "rid-1" { t.Fatalf("OK resp=%+v", r) }

	c2, w2 := newCtx()
	Fail(c2, ErrForb)
	var r2 Response
	json.Unmarshal(w2.Body.Bytes(), &r2)
	if r2.Code != ErrForb || r2.Message == "" { t.Fatalf("Fail resp=%+v", r2) }

	c3, w3 := newCtx()
	FailMsg(c3, ErrBad, "custom")
	var r3 Response
	json.Unmarshal(w3.Body.Bytes(), &r3)
	if r3.Message != "custom" { t.Fatalf("FailMsg resp=%+v", r3) }

	c4, _ := newCtx()
	AbortWithCode(c4, ErrAuth)
	if !c4.IsAborted() { t.Fatal("AbortWithCode did not abort") }
}
