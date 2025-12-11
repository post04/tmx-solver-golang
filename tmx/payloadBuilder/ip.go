package payloadbuilder

import (
	"net/url"

	"github.com/obfio/tmx-solver-golang/tmx"
)

func ip(c *tmx.Client) string {
	// fmt.Println(c.HTTPClient.GetProxy())
	ip, err := c.GetIP()
	if err != nil {
		return ""
	}
	v := &url.Values{}
	v.Set("wei", ip)
	a := "&" + v.Encode()
	return a
}
