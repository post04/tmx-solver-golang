package payloadbuilder

import (
	"net/url"

	"github.com/obfio/tmx-solver-golang/mongo"
)

func ua(p *mongo.Print) string {
	v := &url.Values{}
	v.Set("lq", p.AgentInfo.UserAgent)
	a := v.Encode()
	return a
}
