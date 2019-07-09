package helper

import (
	"strconv"
	"time"
)

func PreNDayTime(n int) int64 {
	return time.Now().AddDate(0, 0, -n).Unix()
}

func ArrayStrToInt(in []string) (outs []int) {
	for _, str := range in {
		if out, err := strconv.Atoi(str); err == nil {
			outs = append(outs, out)
		}
	}
	return
}
