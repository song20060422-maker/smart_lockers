package middleware

import "strconv"

// Name 审计用的操作人标识 @ticket B端-审计
// 定义方法 Name：接收者 u CurrentUser，实现该类型对应的业务行为

func (u CurrentUser) Name() string { return "uid:" + strconv.FormatUint(uint64(u.ID), 10) }
