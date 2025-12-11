package kleinanzeigenMobile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/267H/orderedform"
	"github.com/google/uuid"
	sites_mobile "github.com/obfio/tmx-solver-golang/sites-mobile"
	"github.com/obfio/tmx-solver-golang/tmx/encodePayload"
)

var (
	defaultAndroidPrint = &sites_mobile.AndroidPrint{
		OSVersion: "14",
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
	case 3:
		nonce := guid
		JB := getJB(sessionID, nonce)
		form := orderedform.NewForm(6)
		form.Set("org_id", orgID)
		form.Set("ja", JB)
		form.Set("h", "0")
		form.Set("session_id", sessionID)
		form.Set("nonce", nonce)
		form.Set("m", "2")
		return form.URLEncode()
	default:
		return ""
	}
}

func getJB(sessionID, nonce string) string {
	d := fmt.Sprint(rand.Intn(200) + 800)
	PIDInt := rand.Intn(1000) + 13000
	PID := strconv.Itoa(PIDInt)
	installTime := getInstallTime()
	androidPrint := *defaultAndroidPrint
	form := orderedform.NewForm(73)
	form.Set("hh", makeMD5Hash(orgID+sessionID))
	form.Set("autm", "0")
	// TODO: randomize if needed
	form.Set("aspl", "2023-09-05")
	form.Set("aos", "android")
	form.Set("lq", UserAgent)
	form.Set("cos", "android")
	form.Set("aov", androidPrint.OSVersion)
	form.Set("dm", "false")
	form.Set("pid", PID)
	form.Set("dr", "http://com.ebay.kleinanzeigen")
	form.Set("apd", d)
	form.Set("bhsydpi", "480.0")
	form.Set("bbtm", "0")
	form.Set("sa_pt", firebasePushToken)
	form.Set("ics", "0")
	form.Set("ftsn", fmt.Sprint(rand.Intn(10)+160))
	// TODO: randomize if needed
	form.Set("btst", "{\"level\":1.0,\"status\":\"unplugged\"}")
	form.Set("atr", generateATR(d))
	form.Set("drm", sites_mobile.GenerateDeviceUniqueId())
	form.Set("tzd", "America/New_York")

	form.Set("ats", fmt.Sprint(rand.Intn(999999999)+6000000000))
	form.Set("ab", "samsung")
	form.Set("ad", "e3q")
	form.Set("alo", "en_CA")
	form.Set("mr", "3")
	form.Set("cps", "1,1,1,1")
	// APK hash md5
	form.Set("ah", "6853e65afaba053374729cca199fb57a")
	form.Set("mto", fmt.Sprint(rand.Intn(999999)+3000000))
	form.Set("al", "en-ca")
	form.Set("am", "sm-s928b")
	form.Set("mdf", randomMD5Hash())
	// TODO: randomize if needed
	form.Set("bhsshpx", "2769")
	form.Set("ppid", fmt.Sprint(PIDInt-13000))
	form.Set("upl", "granted:AD_ID,ACCESS_NETWORK_STATE;denied:ACCESS_FINE_LOCATION,ACCESS_COARSE_LOCATION")
	form.Set("at", "agent_mobile")
	form.Set("av", sdkVersion)
	form.Set("apit", fmt.Sprint(installTime))
	form.Set("name", "sdk_gphone64_x86_64")
	form.Set("mds", randomMD5Hash())
	/*
		wc - vpn, NETWORK_INFO_TYPE should probably be WiFi
		ait - 1746924854, STORAGE_EMULATED_TIMESTAMP
		se - enforcing, SELINUX_MODE
		bhsdmo - android-14 samsung:e3q, DEVICE_VERSION_BRAND
		adid - 51a4fbdf-1f61-4a23-9d86-313a8114cd36, ADVERTISING_ID
	*/
	form.Set("amt", fmt.Sprint(installTime+int64(rand.Intn(999999)+2000000)))
	// TODO: might not be random so not sure if this is right
	form.Set("swid", randomMD5Hash())
	form.Set("wc", "wifi")
	form.Set("ait", fmt.Sprint(installTime-int64(rand.Intn(9999)+10000)))
	form.Set("se", "enforcing")
	form.Set("bhsdmo", "android-14 samsung:e3q")
	form.Set("adid", uuid.New().String())
	// TODO: randomize if needed
	form.Set("ipv4", "{\"10.0.2.15\":\"eth0\",\"10.1.10.1\":\"tun0\"}")
	form.Set("pldec1", "{\"description\":\"Not Cloned\"}")
	// TODO: randomize if needed
	form.Set("ipv6", "{\"fe80::5054:ff:fe12:3456\":\"eth0\",\"fe80::7384:f03e:7cb0:db45\":\"tun0\",\"fec0::5054:ff:fe12:3456\":\"eth0\"}")
	// TODO: might be wrong, probably is, I honestly don't know lmao
	form.Set("prst", fmt.Sprint(rand.Intn(999999999999999999)))
	form.Set("ani", randomMD5Hash()[:16])
	// create a function called `generateMEX2` that generates random values for these fields between x and y
	form.Set("mex2", generateMEX2())
	form.Set("c", fmt.Sprint(rand.Intn(999999999999999999)))
	form.Set("mdtm", "0")
	form.Set("fts", randomMD5Hash())
	form.Set("mex6", "64")
	// todo: randomize if needed
	form.Set("f", "2992x1344")
	form.Set("grr", "")
	form.Set("nhc", "4")
	form.Set("anv", "com.simplygood.ct:10.3.1:6100032:x86_64")
	form.Set("ah2", "e97ca36bc0146eb36d65c24372997d997c787ae62e0cde1c295b54d83bc11ba2")
	form.Set("gr", "0")
	form.Set("mrr", "cpu_abi:x86_64;arch:x86_64;prop://init.svc.qemu-props?nil;")
	form.Set("afs", fmt.Sprint(rand.Intn(999999999)+2000000000))
	// make this a timestamp in MS 24 hours ago
	form.Set("abt", fmt.Sprint(time.Now().Add(-24*time.Hour).UnixMilli()))
	form.Set("bhsxdpi", "480.0")
	form.Set("vpn", "false")
	form.Set("w", nonce)
	form.Set("z", "60")
	// TODO: randomize if needed
	form.Set("bhssnby", "6.233333")
	// TODO: randomize if needed
	form.Set("bhssnbx", "2.8")
	form.Set("lh", "http://com.ebay.kleinanzeigen/mobile")
	// TODO: randomize if needed
	form.Set("bhsswpx", "1344")
	a := normalizeEncoding(form.URLEncode())
	// a := form.URLEncode()
	// fmt.Println(a)
	return encodePayload.EncodePayload(a, sessionID)
}

func normalizeEncoding(payload string) string {
	// Replace all uppercase percent encodings with lowercase to match expected format
	replacer := strings.NewReplacer(
		"+", "%20",
		"%2F", "%2f",
		"%3A", "%3a",
		"%2C", "%2c",
		"%3B", "%3b",
		"%7B", "%7b",
		"%7D", "%7d",
		"%22", "%22",
		"%5B", "%5b",
		"%5D", "%5d",
		"%3D", "%3d",
		"%2B", "%2b",
		"%3C", "%3c",
		"%3E", "%3e",
		"%26", "%26",
		"%25", "%25",
		"%23", "%23",
		"%3F", "%3f",
		"%60", "%60",
		"%40", "%40",
		"%24", "%24",
		"%21", "%21",
		"%27", "%27",
		"%28", "%28",
		"%29", "%29",
		"%2A", "%2a",
		"%2D", "%2d",
		"%2E", "%2e",
		"%5C", "%5c",
	)

	return replacer.Replace(payload)
}

// {"mlc":74,"mls":442387908,"slc":247,"sls":85128175,"tda":false}
func generateMEX2() string {
	return fmt.Sprintf("{\"mlc\":%d,\"mls\":%d,\"slc\":%d,\"sls\":%d,\"tda\":%t}",
		rand.Intn(10)+70,
		rand.Intn(99999999)+400000000,
		rand.Intn(20)+230,
		rand.Intn(9999999)+80000000,
		rand.Intn(2) == 0,
	)
}

func getInstallTime() int64 {
	now := time.Now()
	minAgo := now.AddDate(0, 0, -14).UnixMilli()
	maxAgo := now.AddDate(0, 0, -7).UnixMilli()
	return rand.Int63n(maxAgo-minAgo) + minAgo
}

func generateATR(t string) string {
	return fmt.Sprintf("{\"cpo\":%s,\"dyo\":%s,\"psi\":%s,\"pri\":%d,\"cpi\":%s,\"ori\":\"portrait\",\"adb\":0,\"dper\":[\"ACCESS_WIFI_STATE\", \"ACCESS_COARSE_LOCATION\", \"ACCESS_FINE_LOCATION\", \"CHANGE_WIFI_STATE\"],\"mif\":\"\",\"crs\":%d}",
		"9218309884245704190",
		"9218309884245704190",
		"1",
		rand.Intn(100)+300,
		t,
		1,
	)
}

func makeMD5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func randomMD5Hash() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprint(rand.Intn(1000000000000000000))))
	return hex.EncodeToString(h.Sum(nil))
}

/*
hh = md5Hash(org_id + session_id)
autm = 0, 0 means no tampering detected
aspl = 2023-09-05, that is the date of the latest security patch: Build.VERSION.SECURITY_PATCH
aos = android, the operating system of the device
lq = user-agent, Mozilla/5.0 (Linux; Android 14; SM-S928B Build/UE1A.230829.036.A4; wv) AppleWebKit/537.36+ (KHTML, like Gecko) Version/4.0 Chrome/136.0.7103.60 Mobile Safari/537.36+ 7.6-46
cos = android, the operating system of the device
aov = 14, OS version
pid = 13660, agent process ID
dr = http://com.simplygood.ct
apd = 831, time in MS since last request?
bhsydpi = 480.0, DEVICE_DISPLAY_YDPI
bbtm = 0, TAMPER_CODE_BB_MODULE, 0 means no tampering detected
ics = 0, IN_CALL_STATUS, 0 means not in call ?
ftsn = 164, device font count
btst = {"level":1.0,"status":"unplugged"}, DEVICE_BATTERY_STATUS
atr = {"cpo":9218309884245704190,"dyo":9218309884245704190,"psi":1,"pri":343,"cpi":831,"ori":"portrait","adb":1,"dper":["ACCESS_WIFI_STATE", "ACCESS_COARSE_LOCATION", "ACCESS_FINE_LOCATION", "CHANGE_WIFI_STATE"],"mif":"","crs":"0"}
	cpo = 9218309884245704190, option all?
	dyo = 9218309884245704190, option all?
	psi = 1, request number?
	pri = 343, time between something and something else?
	cpi = 831, time in MS since last request?
	ori = portrait, __orientation
	adb = 1, adb detection? I think 0 is good, 1 is bad
	dper = ["ACCESS_WIFI_STATE", "ACCESS_COARSE_LOCATION", "ACCESS_FINE_LOCATION", "CHANGE_WIFI_STATE"], the permissions the app has
	mif = "", hard coded as empty string
	crs = 0, crashlog isn't empty? 1 = is empty

drm - YoYRhd++XQiL8HL7VtVG/XmJC6kF7DTLkkygx6Co3ik=, base64 random 16 bytes = GenerateDeviceUniqueId
tzd - America/New_York, timezone name
ats - 6228115456, total space on the device, probably in bytes
ab - samsung, samsung | phone brand
ad - e3q, phone model
alo - en_CA, locale
mr - 3, em count, something about the CPU
cps - 1,1,1,1 | cpu speed idk how he's generated yet
ah - 6853e65afaba053374729cca199fb57a, APP_SELF_HASH_MD5 not really sure how it's generated but it should just be the same on the same APK
mto - 3050208, memory total in bytes
al - en-ca, language
am - sm-s928b, phone model
mdf - 09e974840db5d0e9e8876cf15a88fe1c, DEVICE_FINGERPRINT not 100% sure what this is but it's probably md5 so it can probably just be a random md5 hash
bhsshpx - 2769, height in pixels
ppid - 374, AGENT_PARENT_PID probably the PID of the actual app instead of the PID of the anti-bot SDK
upl - ["granted:AD_ID,ACCESS_NETWORK_STATE","denied:ACCESS_FINE_LOCATION,ACCESS_COARSE_LOCATION"], checking some permissions, one granted, one denied
at - agent_mobile, agent type
av - 7.6-46, agent version
apit - 1746937569, app install time
name - sdk_gphone64_x86_64, idk but can probably just be static
mds - 36c0670a2325d6dc910cb6df769d6796, DEVICE_STATE not really sure but it's an md5 hash
amt - 1748553128, APP_MODIFICATION_TIME probably the last time the app was updated
swid - 3a98e6e048f94df499e5a30650953670, DEVICE_SOFTWARE_ID not sure what this is but it's an md5 hash
wc - vpn, NETWORK_INFO_TYPE should probably be WiFi
ait - 1746924854, STORAGE_EMULATED_TIMESTAMP
se - enforcing, SELINUX_MODE
bhsdmo - android-14 samsung:e3q, DEVICE_VERSION_BRAND
adid - 51a4fbdf-1f61-4a23-9d86-313a8114cd36, ADVERTISING_ID
ipv4 - {"10.0.2.15":"eth0","10.1.10.1":"tun0"}, internal IP address?
pldec1 - {"description":"Not Cloned"}, PLUGIN_PATH_STR ??? idk man lol
ipv6 - {"fe80::5054:ff:fe12:3456":"eth0","fe80::7384:f03e:7cb0:db45":"tun0","fec0::5054:ff:fe12:3456":"eth0"}, internal IPv6 I guess?
prst - -113494879111127801, TAMPER_CODE_BASE_MODULE not sure what this is but can probably be random?
ani - 0fe3874a32ec5ad0, some sort of ID, maybe android ID?
mex2 - {"mlc":74,"mls":442387908,"slc":247,"sls":85128175,"tda":false}
    mlc - 75, Number of memory-mapped files excluding /system, /dev, and the JNI .so library
	mls - 442387908, represents the total size in bytes of all mapped files that start specifically with the /system directory
	slc - 247, is a counter for the number of memory-mapped files whose paths start with /system
	sls - 85128175, represents the total size in bytes of all mapped files that start specifically with the /system directory
	tda - false, is the .so lobrary loaded? not really sure why it's false, feel like it should be true lol idk
c - -300, timezone offset
mdtm - 0, TAMPER_CODE_DSH_MODULE
fts - 7efbebd3905af29170a29496f77cf103, DEVICE_FONT_HASH
mex6 - 64, DEVICE_ENCRYPTION_STATUS
f - 2992x1344, DEVICE_DISPLAY_RESOLUTION
grr - afile:///system/bin/su;file:///system/bin/su;, ROOT_DETECTION_PATH_STR should probably be empty instead of this
nhc - 4, NUM_OF_CPU_CORES
anv - com.simplygood.ct:10.3.1:6100032:x86_64, AGENT_APP_INFO probably APK name + ":" + version + ":" + build number + ":" + architecture
ah2 - e97ca36bc0146eb36d65c24372997d997c787ae62e0cde1c295b54d83bc11ba2, SHA256 hash of the APK
gr - 2, ROOT_DETECTION_COUNT should probably be 0
mrr - cpu_abi:x86_64;arch:x86_64;prop://init.svc.qemu-props?nil;, EM_PATH_STR not really super sure what this is
afs - 2226126840, DEVICE_FREE_SPACE
abt - 1748548564, DEVICE_BOOT_TIME
bhsxdpi - 480.0, DEVICE_DISPLAY_XDPI
vpn - true, if a VPN is enabled or not, should probably be false
w - 4055faf9c28d47e8, this is the nonce in the XML <N>4055faf9c28d47e8</N>
z - 60, TIMEZONE_DST_DIFF
bhssnby - 6.233333, DEVICE_DISPLAY_NATIVE_BOUND_Y
bhssnbx - 2.8, DEVICE_DISPLAY_NATIVE_BOUND_X
lh - http://com.simplygood.ct/mobile, looks static per app
bhsswpx - 1344, DEVICE_DISPLAY_WIDTH_IN_PIXEL
*/
