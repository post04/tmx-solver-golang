package kleinanzeigenMobileIOS

import (
	"github.com/267H/orderedform"
	sites_mobile "github.com/obfio/tmx-solver-golang/sites-mobile"
)

var (
	defaultAndroidPrint = &sites_mobile.AndroidPrint{
		OSVersion: "18.4.1",
	}
)

func buildPayload(t int, guid, sessionID string) string {
	switch t {
	case 1:
		androidPrint := *defaultAndroidPrint
		form := orderedform.NewForm(6)
		form.Set("os", OS)
		form.Set("osVersion", androidPrint.OSVersion)
		form.Set("thx", guid)
		form.Set("org_id", orgID)
		form.Set("sdk_version", sdkVersion)
		form.Set("session_id", sessionID)
		return form.URLEncode()
	case 2:
		nonce := guid
		form := orderedform.NewForm(4)
		form.Set("org_id", orgID)
		form.Set("session_id", sessionID)
		form.Set("i", "1")
		form.Set("nonce", nonce)
		return form.URLEncode()
	default:
		return ""
	}
}
