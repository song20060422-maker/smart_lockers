package service

import "strconv"

func idstr(id uint) string { return strconv.FormatUint(uint64(id), 10) }
