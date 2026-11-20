package service

import "strconv"

// 定义函数 idstr：封装一段可复用的处理逻辑（参数：id uint）

func idstr(id uint) string { return strconv.FormatUint(uint64(id), 10) }
