package tool

import (
	"time"
)

// ParseTime 将格式为string的时间转为time.Time格式
func ParseTime(timeStr string) (parsedTime time.Time, err error) {
	layouts := []string{
		"2006-01-02 15:04:05",       // 带空格分隔的格式
		"2006-01-02T15:04:05Z",      // ISO 8601 格式
		"2006-01-02T15:04:05Z07:00", // 带时区的 ISO 8601 格式
	}

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, timeStr)
		if err == nil {
			return parsedTime, nil
		}
	}
	return parsedTime, err // 返回最后一个错误
}
