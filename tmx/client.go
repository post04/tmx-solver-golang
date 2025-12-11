package tmx

import (
	"fmt"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Client struct {
	HTTPClient  tlsclient.HttpClient
	Destroy     bool
	Proxy       string
	Auth        string
	Cookies     map[string]string
	URL         *url.URL
	UserAgent   string
	SecCHUa     string
	RetryAmount int
}

func MakeClient(proxy, URL, userAgent string) *Client {
	//proxy = "http://23.230.167.244:3128"
	var tlsProfile profiles.ClientProfile
	tlsProfile = profiles.Chrome_133

	opts := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(60),
		tlsclient.WithInsecureSkipVerify(),
		tlsclient.WithClientProfile(tlsProfile),
		tlsclient.WithNotFollowRedirects(),
	}
	httpClient, err := tlsclient.NewHttpClient(nil, opts...)
	if err != nil {
		panic(err)
	}
	err = httpClient.SetProxy(proxy)
	if err != nil {
		panic(err)
	}
	URL1, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}
	c := &Client{
		HTTPClient: httpClient,
		Proxy:      proxy,
		Cookies:    make(map[string]string),
		URL:        URL1,
		UserAgent:  userAgent,
		SecCHUa:    buildSecCHUA(userAgent),
	}
	return c
}

func (c *Client) MakeHTTPClientHttp1() tlsclient.HttpClient {
	opts := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(60),
		tlsclient.WithInsecureSkipVerify(),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithNotFollowRedirects(),
		tlsclient.WithForceHttp1(),
	}
	httpClient, err := tlsclient.NewHttpClient(nil, opts...)
	if err != nil {
		panic(err)
	}
	err = httpClient.SetProxy(c.Proxy)
	if err != nil {
		panic(err)
	}
	return httpClient
}

type BrandVersion struct {
	Brand   string
	Version string
}

// manually derive brands from parsed agent
func buildSecCHUA(uaString string) string {

	// base logic: assume Chromium-based unless proven otherwise
	var brands []BrandVersion

	// Always add Not:A-Brand
	brands = append(brands, BrandVersion{"Not:A-Brand", "24"})

	if strings.Contains(uaString, "Chrome/") {
		version := parseMajorVersion(uaString, "Chrome/")
		brands = append([]BrandVersion{
			{"Chromium", version},
			{"Google Chrome", version},
		}, brands...)
	} else if strings.Contains(uaString, "Edg/") {
		version := parseMajorVersion(uaString, "Edg/")
		brands = append([]BrandVersion{
			{"Chromium", version},
			{"Microsoft Edge", version},
		}, brands...)
	} else if strings.Contains(uaString, "Brave/") {
		version := parseMajorVersion(uaString, "Brave/")
		brands = append([]BrandVersion{
			{"Chromium", version},
			{"Brave", version},
		}, brands...)
	}

	// format as Sec-CH-UA header
	var parts []string
	for _, b := range brands {
		parts = append(parts, fmt.Sprintf(`"%s";v="%s"`, b.Brand, b.Version))
	}
	return strings.Join(parts, ", ")
}

func parseMajorVersion(ua string, prefix string) string {
	i := strings.Index(ua, prefix)
	if i == -1 {
		return "0"
	}
	i += len(prefix)
	ver := ua[i:]
	dot := strings.Index(ver, ".")
	if dot == -1 {
		return ver
	}
	return ver[:dot]
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
