package encodePayload

import "fmt"

func EncodePayload(payload, key string) string {
	finalStr := fmt.Sprintf("%v&%s", len(payload), payload)
	chars := "0123456789abcdef"
	result := ""
	for i, n := 0, 0; i < len(finalStr); i++ {
		H := charCodeAt(finalStr, i) ^ charCodeAt(key, n)&10
		n++
		if n == len(key) {
			n = 0
		}
		result += charAt(chars, H>>4&15) + charAt(chars, H&15)
	}
	return result
}
