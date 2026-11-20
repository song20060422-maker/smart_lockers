// @ticket B端-脱敏单元测试
package pkg

import "testing"

func TestMaskPhone(t *testing.T) {
	tests := []struct{ in, want string }{
		{"13800138000", "138****8000"},
		{"12345", "***"},
	}
	for _, tt := range tests {
		if got := MaskPhone(tt.in); got != tt.want {
			t.Fatalf("MaskPhone(%q)=%q want %q", tt.in, got, tt.want)
		}
	}
}

func TestMaskNameNo(t *testing.T) {
	if got := MaskName("张三丰"); got != "张**" {
		t.Fatalf("MaskName got %q", got)
	}
	if got := MaskName("李"); got != "李" {
		t.Fatalf("MaskName single got %q", got)
	}
	if got := MaskNo("ABCDEFGH"); got != "AB****GH" {
		t.Fatalf("MaskNo got %q", got)
	}
}
