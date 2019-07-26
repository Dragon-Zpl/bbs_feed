package helper

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
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

// 驼峰转下划线
func Camel2Underline(s string) string {
	re, _ := regexp.Compile("[A-Z]{1}")
	s = re.ReplaceAllStringFunc(s, func(s string) string {
		m := []byte(s)
		return "_" + string(bytes.ToLower(m[0:1]))
	})
	return s
}

// 获取这周开始的时间
func GetWeekStart() time.Time {
	t := time.Now()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	wd := t.Weekday()
	if wd == time.Monday {
		return t
	}
	offset := int(time.Monday - wd)
	if offset > 0 {
		offset -= 7
	}
	return t.AddDate(0, 0, offset)
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

func SuccessWithDataList(datalist interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"dataList": datalist,
		},
	}
}
