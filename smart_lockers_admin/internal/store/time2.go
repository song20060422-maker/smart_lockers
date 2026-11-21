package store

import "time"

// Now 统一时间源，便于 service 调用
// 定义函数 Now：封装一段可复用的处理逻辑（参数：）

func Now() time.Time { return time.Now() }
