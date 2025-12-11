package kleinanzeigenMobileIOS

import (
	"errors"
	"regexp"
	"time"

	uuid "github.com/google/uuid"
	sites_mobile "github.com/obfio/tmx-solver-golang/sites-mobile"
	kleinanzeigenIOS "github.com/obfio/tmx-solver-golang/tmx/mobile/kleinanzeigen-IOS"
)

var nonceRegex = regexp.MustCompile(`<N>[a-z0-9]{16}\<\/N>`)
var requestURLRegex = regexp.MustCompile(`<SPD>[A-z0-9-.]+\<\/SPD>`)
var CSIDRegex = regexp.MustCompile(`<S>[A-Z0-9]{32}</S>`)

func GetCookies(proxy string) (string, error) {
	androidID := makeAndroidID()
	guid := getGUID(androidID)
	sessionID := uuid.New().String()
	client := kleinanzeigenIOS.MakeClient(proxy)
	client.Cookies["thx_guid"] = guid
	client.UserAgent = UserAgent
	encryptedResponse, err := client.MakeRequest(initRequestURL, buildPayload(1, guid, sessionID))
	if err != nil {
		return "", err
	}
	decryptedResponse, err := sites_mobile.DecryptTMXPayload(encryptedResponse, orgID, sessionID)
	if err != nil {
		return "", err
	}
	// get nonce
	nonce := nonceRegex.FindString(decryptedResponse)
	if nonce == "" {
		return "", errors.New("nonce not found")
	}
	nonce = nonce[3 : len(nonce)-4]
	// get CSID
	CSID := CSIDRegex.FindString(decryptedResponse)
	if CSID == "" {
		return "", errors.New("CSID not found")
	}
	CSID = CSID[3 : len(CSID)-4]
	// get request URL
	requestURL := requestURLRegex.FindString(decryptedResponse)
	if requestURL == "" {
		return "", errors.New("requestURL not found")
	}
	requestURL = requestURL[5 : len(requestURL)-6]
	requestURL = "https://" + requestURL + "/fp/clear.png;CIS3SID=" + CSID
	_, err = client.MakeRequest(requestURL, buildPayload(2, nonce, sessionID))
	if err != nil {
		return "", err
	}
	time.Sleep(2 * time.Second)
	client.AttemptLogin(sessionID)
	return guid, nil
}
