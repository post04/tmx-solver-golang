package kleinanzeigen

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) MakeRequest(URL, body string) (string, error) {
	b, err := gzipEncode(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewReader(b))
	if err != nil {
		return "", nil
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Cookie", c.FormatCookies())
	req.Header.Set("Referer", "http://com.ebay.kleinanzeigen")
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate, no-transform")
	req.Header.Set("Accept-Language", "en-CA")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GzipEncode takes an input string and returns its gzip-compressed form.
func gzipEncode(input string) ([]byte, error) {
	var buf bytes.Buffer
	// Create a new gzip.Writer that writes into our buffer
	gw := gzip.NewWriter(&buf)

	// Write the raw bytes to the gzip.Writer
	if _, err := gw.Write([]byte(input)); err != nil {
		gw.Close()
		return nil, err
	}

	// Close to flush all pending data
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Client) AttemptLogin(sessionID string) error {
	body := "{\"loginID\":\"newportsmoker72@yopmail.com\",\"password\":\"Hello1234!\",\"deviceID\":\"aa31eceba4c24adda8cf84627774efb4a177f9d9c5aa9b069a69323aacd8bafd\",\"targetEnv\":\"mobile\",\"remember\":true}"
	req, err := http.NewRequest("POST", "https://m.apim.canadiantire.ca/v1/authorization/signin/rba-tmx", strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CANTIRE/10.2.0 (Android 12; SM-A528B; Samsung SM-A528B; en)")
	req.Header.Set("X-Tmx-Session-Id", sessionID)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	fmt.Println(resp.Status)
	return nil
}
