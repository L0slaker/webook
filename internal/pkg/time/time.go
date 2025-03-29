package time

import (
	"time"
)

const (
	YYYYMMDD = "2006-01-02"
)

func IsValidDate(layout, dateStr string) bool {
	// 使用 time.Parse 验证日期是否合法
	_, err := time.Parse(layout, dateStr)
	return err == nil
}
