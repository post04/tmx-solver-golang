package decodeObfuscation

import (
	"strings"
)

type Decoder struct {
	payloadStr string
	DecodedStr string
}

func decodeSlice(td_I, td_J string) string {
	td_w := []string{""}
	td_O := 0
	for i := 0; i < len(td_J); i++ {
		td_w = append(td_w, fromCharCode(charCodeAt(td_I, td_O)^charCodeAt(td_J, i)))
		td_O++
		if td_O >= len(td_I) {
			td_O = 0
		}
	}
	return strings.Join(td_w, "")
}

func CreateDecoder(payloadStr string) *Decoder {
	td_z := payloadStr[0:32]
	td_m := ""
	test := []int{}
	for i := 32; i < len(payloadStr); i += 2 {
		test = append(test, int(parseInt(payloadStr[i:i+2], 16)))
		td_m += fromCharCode(int(parseInt(payloadStr[i:i+2], 16)))
	}
	return &Decoder{payloadStr: payloadStr, DecodedStr: decodeSlice(td_z, td_m)}
}

func (d *Decoder) Decode(start, end int) string {
	return d.DecodedStr[start : start+end]
}
