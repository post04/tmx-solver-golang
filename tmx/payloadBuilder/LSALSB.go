package payloadbuilder

import "fmt"

func lsalsb(r *Request) string {
	if r.IsLSA {
		a := fmt.Sprintf("lsa=%s", r.LSA)
		return a
	}
	if r.IsLSB {
		a := fmt.Sprintf("lsb=%s", r.LSB)
		return a
	}
	return ""
}
