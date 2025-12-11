package tmx

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (c *Client) MakeRequest(URL string) (string, error) {
	// fmt.Println(URL)
	currURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header = c.FormatHeaders(`Accept: */*|Accept-Encoding: gzip, deflate, br, zstd|Accept-Language: en-US,en;q=0.9|Cache-Control: no-cache|Connection: keep-alive|Host: content.canadiantire.ca|Pragma: no-cache|Referer: https://customerauth.triangle.com/|Sec-Fetch-Dest: script|Sec-Fetch-Mode: no-cors|Sec-Fetch-Site: cross-site|Sec-Fetch-Storage-Access: none|User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36|sec-ch-ua: "Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"|sec-ch-ua-mobile: ?0|sec-ch-ua-platform: "Windows"`)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("sec-ch-ua", c.SecCHUa)
	req.Header.Set("Host", currURL.Host)
	req.Header.Set("Referer", c.URL.Scheme+"://"+c.URL.Host+"/")
	// for name, header := range req.Header {
	// 	fmt.Println(name + ": " + header[0])
	// }
	if len(c.Cookies) != 0 {
		req.Header.Set("Cookie", c.FormatCookies())
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		if strings.Contains(err.Error(), "PROTOCOL_ERROR") {
			c.HTTPClient.CloseIdleConnections()
			c.HTTPClient = c.MakeHTTPClientHttp1()
			return c.MakeRequest(URL)
		}
		return "", err
	}
	defer resp.Body.Close()
	c.SaveCookies(resp.Header)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) SpecialClearPNGRequest(URL, h string) error {
	currURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header = c.FormatHeaders(`Accept-Encoding: gzip, deflate, br, zstd|Accept-Language: en-US,en;q=0.9|Cache-Control: no-cache|Connection: keep-alive|Host: content.canadiantire.ca|Origin: https://customerauth.triangle.com|Pragma: no-cache|Referer: https://customerauth.triangle.com/|Sec-Fetch-Dest: empty|Sec-Fetch-Mode: cors|Sec-Fetch-Site: cross-site|User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36|sec-ch-ua: "Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"|sec-ch-ua-mobile: ?0|sec-ch-ua-platform: "Windows"`)
	req.Header.Set("Accept", h)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("sec-ch-ua", c.SecCHUa)
	req.Header.Set("Host", currURL.Host)
	req.Header.Set("Origin", c.URL.Scheme+"://"+c.URL.Host)
	req.Header.Set("Referer", c.URL.Scheme+"://"+c.URL.Host+"/")
	if len(c.Cookies) != 0 {
		req.Header.Set("Cookie", c.FormatCookies())
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.SaveCookies(resp.Header)
	return nil
}

func (c *Client) GetIP() (string, error) {
	req, err := http.NewRequest("GET", "https://api.ipify.org/", nil)
	if err != nil {
		return "", err
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// c.Destroy = true
	// fmt.Println("DEBUG: IP: " + string(b) + " == " + c.Proxy)
	return string(b), nil
}

var letterRunes = []rune("abcdef0123456789")

func ctRandomString() string {
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "TRI_" + string(b)
}

var letterRunes1 = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randStrLen(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes1[rand.Intn(len(letterRunes1))]
	}
	return string(b)
}

func getClientUUID() string {
	a := strings.Split("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx", "")
	for i, p := range a {
		if p != "x" && p != "y" {
			continue
		}
		r := int64(rand.Float64()*16) | 0
		t := int64(0)
		if p == "x" {
			t = r
		} else {
			t = r&0x3 | 0x8
		}
		a[i] = strconv.FormatInt(t, 16)
	}
	return strings.Join(a, "")
}

func (c *Client) GenerateRandomSessionID(siteKey string) string {
	switch siteKey {
	case "canadiantire":
		return ctRandomString()
	case "citi":
		return fmt.Sprintf(`%s`, randStrLen(64))
	case "gyft":
		return randStrLen(32)
	case "okx":
		return fmt.Sprintf(`%s_%v`, randStrLen(11), time.Now().UnixMilli())
	case "skrill":
		return getClientUUID()
	case "vanillagift":
		return getClientUUID()

	default:
		return ""
	}
}

func (c *Client) CTTestLogin(sessionID string) bool {
	payload := fmt.Sprintf(`{"loginID":"%s@gmail.com","password":"Lolxd123!","remember":true,"targetEnv":"browser"}`, randStrLen(11))
	req, err := http.NewRequest("POST", "https://apim.triangle.com/v1/authorization/signin/rba-tmx", strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Bannerid", "TRIANGLE")
	req.Header.Set("Basesiteid", "TRIANGLE")
	req.Header.Set("Browse-Mode", "undefined")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ocp-apim-subscription-key", "dfa09e5aac1340d49667cbdaae9bbb3b")
	req.Header.Set("origin", "https://customerauth.triangle.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("referer", "https://customerauth.triangle.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 OPR/120.0.0.0")
	req.Header.Set("sec-ch-ua", `"Opera GX";v="120", "Not-A.Brand";v="8", "Chromium";v="135"`)
	req.Header.Set("service-client", "ctr/web")
	req.Header.Set("service-version", "ctc-dev2")
	req.Header.Set("x-tmx-session-id", sessionID)
	req.Header.Set("x-web-host", "customerauth.triangle.com")
	req.Header.Set("Cookie", "_gid=GA1.2.1616897038.1759513782; kndctr_A6C5776A5245B09C0A490D44_AdobeOrg_cluster=or2; kndctr_A6C5776A5245B09C0A490D44_AdobeOrg_identity=CiYzOTA0MzIxMDcyNDQyMTMxNDU5MTM4NzI3MzgwMzY5NTAzNDg4OFISCJezxtmaMxABGAEqA09SMjAA8AGXs8bZmjM=; AMCV_A6C5776A5245B09C0A490D44%40AdobeOrg=MCMID|39043210724421314591387273803695034888; _scid=iuPQsgA7CrSspB_UPfd6aJj4nDO4qvTy; gig_bootstrap_3_HJ80G1HPNNiHj1XllefNdB0vY4Dg9m66z8LMwITS6NpDTmD8OVJXE3GyeGjCaoM6=secure-gigya_ver4; QuantumMetricSessionID=202a7eb98a61c1ff46314d0c09943cdf; QuantumMetricUserID=33c1f60ebbffe80869dfa9d3d96d1aea; bm_mi=363D14A353C0674D0423BD8C77EE0A06~YAAQBmnMFzQ935eZAQAApIgyqx21StmObT7oXDnf+4cauWHwzcvs2Yqn8kiogHgjyXSzzZBv+2SUIIUXcd8ESBAeO2zV0NcrJBoDGyarjxPllWOu5LwqnlsMYXXK9SUIqDIKZG0WKfa39nidr0PLjYpvZ0PrjR6fOeh1OvMN6SYfJLfdajhcQIXt4e8CuOQYE2nNd2XzP1AZXIYS28m2aolxqmGCbObPByXdplQ+6zhRqQ17mMsGxcKur6WIKnDayCs1OldWHj+Bky9ao+62WWAxfmKGOxVF6RLT7qmzDUBQSi8cY/hF98dEEa8B7uk33q8eiOX7LrdEB3bGp8l087spGHJk4//LvC4=~1; _abck=F7AE9B1AE676CE0E6D2C64514CCC6E3F~0~YAAQBmnMFzk935eZAQAALYkyqw7rS07wVNYFH60RXNojIt/c0jG/u1Vq0hOrf56JBmPzu8iWww/4x/jlQwVWmyk3v8BRXrxayTeNrTIbEhnfy+w8cqttsGKWKEUBr8JSyZdXSwF3ipOFp4PmGWpEmwVyV9hxMzXmfwUiZ+tTc0qBXtg2bXnH0cUx8fh5FI6VyrvDjP6cieBZAxfwnohqrxJpvSvWPhtNd7UwvylA4HSmApGUkDsdM3/rx2Jm+FpGuYPkoNmAsVyGpNzB3xOmck/3/0M5niddlnZDzKoWCMvCYJgR9p2caojTlpGDPLe5gBMSN2TG88x4hqBRePKSr69hun4zOHSDlYdglX+z1N3A2BRgOEcu6yikco9lelALKa4yswri58RME3ACHnU08XaUrrBQfGQeTBReYfMumWt6ycNA3tZufj2B/7Z+X04rxPjD8XMligxaxIJ2RCUGzX9rwqftP9rHCdfwARb4CvfPXJ32e0BxaCyRR9r5JDZKBBIM1ro4ImqznwOJGu1A4IK06hDTS3IPYBdP5SxNMHAeTXE7dN+JQbwfsYid59Dg/KfhtxAp/ayAjV159auPsvU0OV8UIrsIH6OnDM4lXYik/M2QxAJJNDqm/cNzmSapEo7BzelKwEyMJhW1d6oif3YoRYbzMJeiv/wNjUTK~-1~-1~1759517378~AAQAAAAE%2f%2f%2f%2f%2f7TREFm8ajS9qoVFbw1vHsuf7p8WE3m79aeBW+ZIyexC6qAwWxLaCkMeFUl0ASyQ3k8tAQLq8bFsdFq5+40ZRslSawH3e1eSnDOJ57npfwoWlkxiTd1YeYCqt8tPb1wHtlZOfXM%3d~-1; _gat_UA-12123123-123=1; ak_bmsc=D578B7FA19D1527A8CEDEADAF86E6C81~000000000000000000000000000000~YAAQBmnMF5I935eZAQAAnYwyqx06kSYUrOekCwOBEFUO2v2nTJUe8C5wEhrQIPAP6uhucdRQUQnpohb/SJI/iU3QS1c8SIvwOYHsl3IDrhqIp2lGxxBkofpD1jvlGdmc2JJaYIs8VK5At1muD0JInBu99A05UODJFbMZO5PaRet//lUhl+kpogKeGDdnHs99W/UhLQsRtkPkzzQRtjxi5xsdNeEQoSkF1ZPc76CV5JahjuAVowmM/7XO/A0GqWOkRCdg22Td7aA9JQ7IUUb5EK8RTZzz85eJLHsLmOE0NsVqc95UxS2IIdbqW7sSD0Mty0wgDwA8aG8BXh1o0xSDZZ2orbBicKUimNzyhaf0uU2FIhi4ec90koKFTD1rPBq2Lbd51pKI+68QoNedv1pNHHY6b/wF4odkYa92XdjFGyeWoH7cMio8jZ5uHdrVDSFYErcJD9t2YoPlzSWzqYHJIfe5uOWZuJNSNaZT5F29nTPiOMD75vRO/dRRyPZ9xymMUKleHmsBpNmZsNapZo37P/0IzjGDxJIDmePKQbYOnIaJlFc3tjc=; gig_bootstrap_3_6zHlynZz9NnduhqUerMynPILRICcOo0i_6JBR7PLXEsxfFH0jSB0dIFoj1XM9CIy=secure-gigya_ver4; bm_sz=7A8286B8D7DC5D8E3B97AF8D4ED793B9~YAAQBmnMF3k/35eZAQAAfJ8yqx3uefrDqMD8/t5qzOcOHBbuCwbrI5CfI7mq58Zb/pcqP0DnbVtwJka8x4jba+DlHznbKZTDiEiAHI+Y4O4Hg3vFwYzTI1de0SbcmRSlmOB9YOozQ+ooAPL7v6TxFX2dqvaCjGHsy8gtyq2yafS1T+ejFyAu+pTW+IpgxlHVBOZSNQVfRgCmW3vfqFnlW8plrlUihKYNNJinohhB2AKOV0+Al3yOEsFY9f+nbvpfYkHrg11iymfQUyeiHJSMj3PFhUulgaP8caPQGw9Eh9GY0Qi3Cy/V/kupc4lZMQLOeU0lsCymxEJajkYx/OCdU30UfXfw9pWCKW/QjBRGfjV6sid+SU21TQC9FoF9a/OTVPaiZBdJkahSevohb8/toVKiUTeaCFYGapr6y4/f1IWzRyt4gHz0ilaQOnTV8bHU4MCPP9usVlaeQDiS7/9M~4539959~4277552; _scid_r=lGPQsgA7CrSspB_UPfd6aJj4nDO4qvTy-O-HqA; _ga=GA1.1.113939590.1759513782; _ga_EXYMG8B1M5=GS2.1.s1759513782$o1$g1$t1759513850$j53$l0$h0; bm_sv=4688273EFBF8B846A7FDD277BA9AAFD9~YAAQBmnMF5pA35eZAQAAtasyqx3DdT37GYs/sNXPKpbI8des0KSdgDt38tb5rr+LGeVFHDqNMwkMTvPibbA22IieuFvBI8ZTLFG6D4BG0/mfS3cCWimMpJnXLm1L55+3QZ9c+8eWWgD8H9nH2b9TYozvnTKA+1q9g9kzzH3yeQRk83LplzO3sISei8xJtHxYIxeZbYUjwWdRASXJSUSl6K+szI343lIlQkUb8Khzj52xX67HkCCDGAawVGnhhI0qNGGN~1; _gcl_au=1.1.106823120.1759513782.930111867.1759513789.1759513855")
	resp, err := c.DoRequest(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b) + c.Proxy)
	if strings.Contains(string(b), "invalid loginID or password") {
		fmt.Println("DEBUG: Proxy used is good: " + c.Proxy)
		return true
	}
	return false
}

func (c *Client) CTTestLoginMobile(sessionID string) bool {
	//https://apim.canadiantire.ca/v1/authorization/signin/rba-tmx
	//{\"remember\":true,\"targetEnv\":\"mobile\",\"loginID\":\"<input.USER>\",\"password\":\"<input.PASS>\",\"deviceID\":\"<device>\"}
	payload := fmt.Sprintf(`{"remember":true,"targetEnv":"mobile","loginID":"%s@gmail.com","password":"Lolxd123!","deviceID":"%s"}`, randStrLen(11), randStrLen(32))
	req, err := http.NewRequest("POST", "https://apim.canadiantire.ca/v1/authorization/signin/rba-tmx", strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Tmx-Session-Id", sessionID)
	resp, err := c.DoRequest(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b) + c.Proxy)
	return strings.Contains(string(b), "invalid loginID or password")
}
