package utils

import "time"

func AbsoluteTimeFromNow(duration time.Duration) int64 {
	return time.Now().Add(duration).Unix()
}
