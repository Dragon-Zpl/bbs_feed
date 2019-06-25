package stringi

import (
	"strconv"
	"strings"
)

type Form map[string]string

// 字符串模板
func Build(message string, bind Form) string {
	for k, v := range bind {
		var re = "{" + k + "}"
		message = strings.Replace(message, re, v, -1)
	}
	return message
}

func Swap(a string, b string) (string, string) {
	return b, a
}

func Reverse(arr []string) {
	var n int
	var length = len(arr)
	n = length / 2
	for i := 0; i < n; i++ {
		arr[i], arr[length-1-i] = Swap(arr[i], arr[length-1-i])
	}
}

// 转义引号
func AddSlashes(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "`", "\\`", -1)
	return str
}

func Empty(str string) bool {
	str = strings.TrimSpace(str)
	return (str == "") || (str == "0")
}

func ToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

func ToInt64(s string) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
