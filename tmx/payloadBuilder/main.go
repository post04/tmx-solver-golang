package payloadbuilder

import (
	"github.com/obfio/tmx-solver-golang/mongo"
	"github.com/obfio/tmx-solver-golang/tmx"
)

type RequestType int

const (
	RequestTypeBrowserVer RequestType = iota
	RequestTypeBrowserGeneral
	RequestTypeUA
	RequestTypeIP
	RequestTypeVideoAudio
	RequestTypeSID
	RequestTypeWGL
	RequestTypeLSALSB
)

type Request struct {
	Print     *mongo.Print
	SessionID string
	URL       string
	Nonce     string
	T         RequestType
	OrgID     string
	Client    *tmx.Client
	LSA       string
	LSB       string
	IsLSA     bool
	IsLSB     bool
}

func BuildPayload(r *Request) string {
	switch r.T {
	case RequestTypeBrowserVer:
		return browserVer(r.Print)
	case RequestTypeBrowserGeneral:
		return browserGeneral(r)
	case RequestTypeUA:
		return ua(r.Print)
	case RequestTypeIP:
		return ip(r.Client)
	case RequestTypeVideoAudio:
		return videoAudio(r.Print)
	case RequestTypeSID:
		return sid(r)
	case RequestTypeWGL:
		return wgl(r.Print)
	case RequestTypeLSALSB:
		return lsalsb(r)
	}
	return ""
}
