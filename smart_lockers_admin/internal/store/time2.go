package store

import "time"

// Now 统一时间源，便于 service 调用
func Now() time.Time { return time.Now() }
