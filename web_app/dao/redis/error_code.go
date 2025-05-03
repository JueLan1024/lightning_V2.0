package redis

import "errors"

var (
	ErrorRefreshTokenNotFound = errors.New("RefreshToken未找到")
	ErrorInvalidDataFormat    = errors.New("获取的数据格式不正确")
	ErrorParseDataFailed      = errors.New("解析数据失败")
	ErrorDataNotFound         = errors.New("未找到数据")
)
