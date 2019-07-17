package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func PreNDayTime(n int) int64 {
	return time.Now().AddDate(0, 0, -n).Unix()
}

func PreMinuteTime(duration time.Duration) int64 {
	return time.Now().Add(-duration).Unix()
}

func ArrayStrToInt(in []string) (outs []int) {
	for _, str := range in {
		if out, err := strconv.Atoi(str); err == nil {
			outs = append(outs, out)
		}
	}
	return
}

func SameElementCount(s []int) map[int]int {
	m := make(map[int]int)
	for _, elem := range s {
		if _, ok := m[elem]; !ok {
			m[elem] = 1
		} else {
			m[elem] += 1
		}
	}
	return m
}

func Success() (int, interface{}) {
	return http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	}
}

func SuccessWithDate(data interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	}
}
