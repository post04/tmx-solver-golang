package tmx

import (
	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.RetryAmount++
		if c.RetryAmount < 4 {
			return c.DoRequest(req)
		}
		return nil, err
	}
	return resp, nil
}
