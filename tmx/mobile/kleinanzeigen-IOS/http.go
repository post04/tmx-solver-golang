package kleinanzeigenIOS

import (
	"bytes"
	"compress/gzip"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io/ioutil"
	"strings"
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
	req.Header.Set("Referer", "http://Kleinanzeigen/")
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
	body := "{\"username\":\"eveleen49@dcpa.net\",\"password\":\"Lolxd123!!!!\"}"
	req, err := http.NewRequest("POST", "https://gateway.kleinanzeigen.de/auth/login", strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Kleinanzeigen/100.47.0 (com.ebaykleinanzeigen.ebc; build:25.136.10553473; iOS 18.4.1) Alamofire/5.10.2")
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
