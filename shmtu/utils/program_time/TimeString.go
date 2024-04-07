package program_time

import (
	"shmtu-cas-go/shmtu/utils/str"
	"strings"
	"time"
)

func AddChatTo6DigitTime(timeStr string) string {
	parts := str.SplitLengthN(timeStr, 2)
	result := strings.Join(parts, ":")
	return result
}

func ParseTimeFromString(dateTimeStr string) (time.Time, error) {
	return time.Parse("2006.01.02 15:04:05", dateTimeStr)
}
