package kleinanzeigenIOS

import (
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Client struct {
	HTTPClient tlsclient.HttpClient
	Proxy      string
	Auth       string
	Cookies    map[string]string
	Site       string
	UserAgent  string
}

func MakeClient(proxy string) *Client {
	//proxy = "http://23.230.167.244:3128"
	var tlsProfile profiles.ClientProfile
	tlsProfile = profiles.ConfirmedIos

	opts := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(60),
		tlsclient.WithInsecureSkipVerify(),
		tlsclient.WithClientProfile(tlsProfile),
		tlsclient.WithNotFollowRedirects(),
		tlsclient.WithForceHttp1(),
	}
	httpClient, err := tlsclient.NewHttpClient(nil, opts...)
	if err != nil {
		panic(err)
	}
	err = httpClient.SetProxy(proxy)
	if err != nil {
		panic(err)
	}
	c := &Client{
		HTTPClient: httpClient,
		Proxy:      proxy,
		Cookies:    make(map[string]string),
	}
	return c
}

// FormatHeaders turns a string of headers seperated by `|` into a http.Header map
func (c *Client) FormatHeaders(h string) http.Header {
	headers := http.Header{}
	for _, header := range strings.Split(h, "|") {
		parts := strings.Split(header, ": ")
		headers.Set(parts[0], parts[1])
	}
	return headers
}

// FormatCookies takes the cookies currently stored in s.Cookies map and turns them into a string value
func (c *Client) FormatCookies() string {
	cookies := ""
	for key, value := range c.Cookies {
		cookies += key + "=" + value + "; "
	}
	if len(cookies) > 3 {
		cookies = cookies[:len(cookies)-2]
	}
	return cookies
}

// SaveCookies saves the cookies from a current request into the s.Cookies map
func (c *Client) SaveCookies(cookies http.Header) {
	for _, cookie := range cookies["Set-Cookie"] {
		parts := strings.Split(cookie, "; ")
		c.Cookies[strings.Split(parts[0], "=")[0]] = strings.Join(strings.Split(parts[0], "=")[1:], "=")
	}
}
