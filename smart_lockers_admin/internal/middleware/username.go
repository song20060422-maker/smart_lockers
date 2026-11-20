package middleware

import "strconv"

// Name 审计用的操作人标识 @ticket B端-审计
func (u CurrentUser) Name() string { return "uid:" + strconv.FormatUint(uint64(u.ID), 10) }
