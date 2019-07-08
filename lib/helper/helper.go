package helper

import "time"

func PreNDayTime(n int) int64 {
	return time.Now().AddDate(0, 0, -n).Unix()
}
