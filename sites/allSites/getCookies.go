package allsites

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/267H/orderedform"
	"github.com/mileusna/useragent"
	geturls "github.com/obfio/tmx-solver-golang/getURLs"
	"github.com/obfio/tmx-solver-golang/mongo"
	"github.com/obfio/tmx-solver-golang/sites"
	"github.com/obfio/tmx-solver-golang/tmx"
	"github.com/obfio/tmx-solver-golang/tmx/encodePayload"
	payloadbuilder "github.com/obfio/tmx-solver-golang/tmx/payloadBuilder"
)

var (
	Sites map[string]Site
)

type Site struct {
	Name                 string `json:"name"`
	OrgID                string `json:"orgID"`
	PageID               string `json:"pageID"`
	URL                  string `json:"url"`
	RandomInit           bool   `json:"randomInit"`
	GenerateOwnSessionID bool   `json:"generateOwnSessionID"`
	HP                   struct {
		Enabled      bool   `json:"enabled"`
		Regex        string `json:"regex"`
		RegexUseable *regexp.Regexp
	} `json:"HP"`
	ARF struct {
		Enabled  bool     `json:"enabled"`
		Elements []string `json:"elements"`
	} `json:"ARF"`
}

func init() {
	// read the file sites.json
	f, err := os.ReadFile("./sites.json")
	if err != nil {
		log.Fatal(err)
	}

	var sites []Site
	err = json.Unmarshal(f, &sites)
	if err != nil {
		log.Fatal(err)
	}

	Sites = make(map[string]Site, len(sites))
	for _, site := range sites {
		if site.HP.Enabled {
			site.HP.RegexUseable = regexp.MustCompile(site.HP.Regex)
		}
		Sites[site.Name] = site

	}
}

type Response struct {
	SessionID string `json:"sessionID"`
	UserAgent string `json:"userAgent"`
	Proxy     string `json:"proxy"`
	Error     string `json:"error"`
}

// Request and response structures for the API
type ProxyRequest struct {
	APIKey      string `json:"apiKey"`
	Proxy       string `json:"proxy"`
	Site        string `json:"site"`
	SessionID   string `json:"uuid"`
	URL         string `json:"url"`
	Mobile      bool   `json:"mobile,omitempty"`
	StopAt      int    `json:"stopAt,omitempty"`
	ARFDisabled bool   `json:"arfDisabled,omitempty"`
}

var (
	MaxStopAt = 16
)

func GetCookies(proxyRequest *ProxyRequest) *Response {
	proxyRequest.Proxy = strings.ReplaceAll(proxyRequest.Proxy, "{rand}", RandStringBytesMaskImprSrc1(16))
	if proxyRequest.StopAt < 1 || proxyRequest.StopAt > MaxStopAt {
		proxyRequest.StopAt = MaxStopAt
	}
	response := &Response{Proxy: proxyRequest.Proxy}

	site := Sites[proxyRequest.Site]

	st := time.Now()
	errors := atomic.Int64{}
	done := atomic.Int64{}
	if proxyRequest.ARFDisabled {
		done.Add(1)
	}

	print := mongo.GetRandomPrint()
	response.UserAgent = print.AgentInfo.UserAgent

	builderRequest := &payloadbuilder.Request{
		Print:     print,
		SessionID: proxyRequest.SessionID,
		URL:       proxyRequest.URL,
		Nonce:     "",
		T:         payloadbuilder.RequestTypeBrowserVer,
		OrgID:     site.OrgID,
	}

	client := tmx.MakeClient(proxyRequest.Proxy, proxyRequest.URL, print.AgentInfo.UserAgent)
	builderRequest.Client = client
	if site.GenerateOwnSessionID && proxyRequest.SessionID == "" {
		builderRequest.SessionID = client.GenerateRandomSessionID(site.Name)
	}

	response.SessionID = builderRequest.SessionID

	// tags.js
	initType := "init"
	if site.RandomInit {
		initType = "randomInit"
	}
	initURL := buildURL(initType, builderRequest.SessionID, proxyRequest.URL, []string{}, site)
	// fmt.Println(initURL + "initURL")
	initScript, err := client.MakeRequest(initURL)
	if err != nil {
		response.Error = "Failed to get cookies: " + err.Error()
		return response
	}
	firstScriptURLs, err := geturls.GetTagDynamic(initScript)
	if err != nil {
		response.Error = "Failed to get cookies: " + err.Error()
		return response
	}
	builderRequest.SessionID = strings.ToLower(builderRequest.SessionID)

	// check.js, should have JB
	builderRequest.T = payloadbuilder.RequestTypeBrowserVer
	payload := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	mainScript, err := client.MakeRequest(buildURL("basicJB", "", firstScriptURLs.CheckJSURL, []string{payload}, site))
	if err != nil {
		response.Error = "Failed to get cookies: " + err.Error()
		return response
	}
	secondScriptURLs, err := geturls.GetCheckJSDynamic(mainScript)
	if err != nil {
		response.Error = "Failed to get cookies: " + err.Error()
		return response
	}

	// m = 2
	go doRequestConcurrently(client, firstScriptURLs.ClearPNGURL, &errors, &done)
	// m = 1
	go doRequestConcurrently(client, firstScriptURLs.ClearPNG2URL, &errors, &done)

	// HP request if HP is enabled
	HPExtractedURL := ""
	if site.HP.Enabled {
		HPHTML, err := client.MakeRequest(secondScriptURLs.HPURL)
		if err != nil {
			response.Error = "Failed to get cookies: " + err.Error()
			return response
		}
		checkJSHPURL := site.HP.RegexUseable.FindString(HPHTML)
		if checkJSHPURL == "" {
			response.Error = "Failed to get cookies: checkJSHPURL == \"\""
			return response
		}
		HPExtractedURL = checkJSHPURL
	}

	// clear.png with special accept header
	go func() {
		err = client.SpecialClearPNGRequest(secondScriptURLs.BlankClear.URL, fmt.Sprintf(`*/*, %s/%s%s`, site.OrgID, secondScriptURLs.BlankClear.Nonce, builderRequest.SessionID))
		if err != nil {
			errors.Add(1)
			return
		}
		done.Add(1)
	}()

	// ls_fp.html
	finalScript, err := client.MakeRequest(secondScriptURLs.LSFPURL)
	if err != nil {
		response.Error = "Failed to get cookies: " + err.Error()
		return response
	}

	// LSA JB
	go doRequestConcurrently(client, secondScriptURLs.LSA.URL+"&jb="+encodePayload.EncodePayload(payloadbuilder.BuildPayload(&payloadbuilder.Request{IsLSA: true, LSA: strings.Split(secondScriptURLs.LSA.Nonce, "_")[0], T: payloadbuilder.RequestTypeLSALSB}), builderRequest.SessionID), &errors, &done)

	// only do this is sidfpURL is not "FAILED"
	if secondScriptURLs.SidFPURL != "FAILED" {
		go func() {
			sidHTML, err := client.MakeRequest(secondScriptURLs.SidFPURL)
			if err != nil {
				errors.Add(1)
				return
			}
			sidDynamic, err := geturls.GetSidFPDynamic(sidHTML)
			if err != nil {
				errors.Add(1)
				return
			}
			// SID from sid_fp.html
			builderRequest.T = payloadbuilder.RequestTypeSID
			builderRequest.Nonce = sidDynamic.SID.Nonce
			sidBody := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
			// fmt.Println(sidDynamic.SID.URL, "AAAAAAAAAAAH")
			doRequestConcurrently(client, buildURL("SID", builderRequest.SessionID, sidDynamic.SID.URL, []string{sidBody}, site), &errors, &done)
		}()
	} else {
		errors.Add(1)
	}

	thirdScriptURLs, err := geturls.GetLsFpDynamic(finalScript)
	if err != nil {
		return response
	}
	// es.js
	go doRequestConcurrently(client, secondScriptURLs.ESJSURL, &errors, &done)

	// top_fp.html
	go doRequestConcurrently(client, secondScriptURLs.TOPFPHTMLURL, &errors, &done)

	// JA JB
	builderRequest.T = payloadbuilder.RequestTypeBrowserGeneral
	jaPayload := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	builderRequest.T = payloadbuilder.RequestTypeUA
	jbPayload := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	go doRequestConcurrently(client, buildURL("JBJA", builderRequest.SessionID, secondScriptURLs.JAJBURL, []string{jaPayload, jbPayload}, site), &errors, &done)

	// if ARF is enabled, do ARF
	ARFScript := ""
	if site.ARF.Enabled && !proxyRequest.ARFDisabled {
		ARFScript, err = client.MakeRequest(HPExtractedURL)
		if err != nil {
			return response
		}
	}

	// h64 request
	go doRequestConcurrently(client, secondScriptURLs.H64URL, &errors, &done)
	// LSB
	go doRequestConcurrently(client, thirdScriptURLs.LSB.URL+"&jf="+encodePayload.EncodePayload(payloadbuilder.BuildPayload(&payloadbuilder.Request{IsLSB: true, LSB: strings.Split(thirdScriptURLs.LSB.Nonce, "_")[0], T: payloadbuilder.RequestTypeLSALSB}), builderRequest.SessionID), &errors, &done)
	// ES.js with empty fr body
	go doRequestConcurrently(client, secondScriptURLs.ESJSURL+"&fr=", &errors, &done)
	// long subdomain
	go doRequestConcurrently(client, secondScriptURLs.LongSubdomainURL, &errors, &done)

	// clear3.png audio video stuff
	builderRequest.T = payloadbuilder.RequestTypeVideoAudio
	videoAudioPayload := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	go doRequestConcurrently(client, buildURL("JE", builderRequest.SessionID, secondScriptURLs.VIDEOAUDIOURL, []string{videoAudioPayload}, site), &errors, &done)
	// SID clear1.png
	builderRequest.T = payloadbuilder.RequestTypeSID
	builderRequest.Nonce = secondScriptURLs.SID.Nonce
	sidBody := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	go doRequestConcurrently(client, buildURL("SID", builderRequest.SessionID, secondScriptURLs.SID.URL, []string{sidBody}, site), &errors, &done)
	// WGL clear.png
	builderRequest.T = payloadbuilder.RequestTypeWGL
	wglBody := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	go doRequestConcurrently(client, buildURL("WGL", builderRequest.SessionID, secondScriptURLs.WGLURL, []string{wglBody}, site), &errors, &done)

	if site.ARF.Enabled && !proxyRequest.ARFDisabled {
		arfURLs, err := geturls.GetARFDynamic(ARFScript)
		if err != nil {
			response.Error = "Failed to get cookies: " + err.Error()
			return response
		}
		arfURL := arfURLs.ARFURL
		form := orderedform.NewForm(3)
		uaInfo := useragent.Parse(print.AgentInfo.UserAgent)
		defaultProfile := &sites.UserAgent{
			BrowserName:    uaInfo.Name,
			BrowserVersion: fmt.Sprint(uaInfo.VersionNo.Major),
			OSName:         uaInfo.OS,
			OSVersion:      uaInfo.OSVersion,
		}
		form.Set("sera_parametere", sites.MakeARFPayload(ARFScript, arfURLs.ARFNonce, builderRequest.SessionID, defaultProfile, site.ARF.Elements))
		form.Set("count", "0")
		form.Set("max", "0")
		arfOutput, err := client.MakeRequest(arfURL + "&" + form.URLEncode())
		if err != nil {
			response.Error = "Failed to get cookies: " + err.Error()
			return response
		}
		if !strings.Contains(arfOutput, "authentic site") {
			fmt.Println(arfOutput, "arfOutput - "+site.Name)
		}
	}

	// IP payload, final request
	builderRequest.T = payloadbuilder.RequestTypeIP
	IPPayload := encodePayload.EncodePayload(payloadbuilder.BuildPayload(builderRequest), builderRequest.SessionID)
	go doRequestConcurrently(client, secondScriptURLs.IPURL+"&jac=1&je="+IPPayload, &errors, &done)
	// finally, wait for all requests to finish
	for {
		if time.Since(st) > 60*time.Second {
			client.HTTPClient.CloseIdleConnections()
			response.Error = "Failed to get cookies: Timeout > 60s"
			return response
		}
		if errors.Load()+done.Load() >= int64(proxyRequest.StopAt) {
			// client.CTTestLogin(response.SessionID)
			client.HTTPClient.CloseIdleConnections()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("Finished in %s (%s): %s\n", time.Since(st), site.Name, "REDACTED")
	return response
}

func buildURL(t, sessionID, URL string, body []string, site Site) string {
	switch t {
	case "randomInit":
		return getInitURL(site.OrgID, site.PageID, site.URL, sessionID)
	case "init":
		return fmt.Sprintf(`%stags.js?org_id=%s&session_id=%s`, site.URL, site.OrgID, sessionID)
	case "basicJB":
		return fmt.Sprintf(`%s&jb=%s`, URL, body[0])
	case "JBJA":
		return fmt.Sprintf("%s&ja=%s&jb=%s", URL, body[0], body[1])
	case "JE":
		return fmt.Sprintf("%s&bbv=3&jac=1&je=%s", URL, body[0])
	case "SID":
		return fmt.Sprintf("%s&jf=%s", URL, body[0])
	case "WGL":
		return fmt.Sprintf("%s&jac=1&je=%s", URL, body[0])
	}

	return ""
}

func doRequestConcurrently(client *tmx.Client, URL string, errors *atomic.Int64, done *atomic.Int64) {
	_, err := client.MakeRequest(URL)
	if err != nil {
		errors.Add(1)
		return
	}
	done.Add(1)
}

const letterBytes = "abcdef0123456789"
const (
	letterIdxBits = 4                    // 4 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src1 = rand.NewSource(time.Now().UnixNano())

// RandStringBytesMaskImprSrc returns a random hexadecimal string of length n.
func RandStringBytesMaskImprSrc1(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src1.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src1.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
