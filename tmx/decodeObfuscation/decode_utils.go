package decodeObfuscation

import (
	"strconv"
)

// CharCodeAt https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/charCodeAt
func charCodeAt(str string, n int) int {
	if len(str) == 0 || len(str) < n {
		return 0
	}
	return int([]rune(str)[n])
}

// FromCharCode https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/FromCharCode
func fromCharCode(c int) string {
	return string(rune(c))
}

func toInt(v interface{}) int {
	switch c := v.(type) {
	case int:
		return int(c)
	case int32:
		return int(c)
	case int64:
		return int(c)
	case float32:
		return int(c)
	case float64:
		return int(c)
	case string:
		r, _ := strconv.Atoi(c)
		return r
	default:
		return 0
	}
}

func parseInt(v, c interface{}) int64 {
	solved, _ := strconv.ParseInt(v.(string), toInt(c), 64)
	return solved
}
