package sites

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"

	"github.com/267H/orderedform"
)

func MakeARFPayload(script, nonce, sessionID string, agent *UserAgent, htmlElements []string) string {
	form := orderedform.NewForm(14)
	form.Set("fp", getFPHash(htmlElements, nonce))
	form.Set("ss", "inputs={} jselements={0} hosts={}")
	form.Set("di", "495238ca2c6d26fc643990c07e313281a0bce6cc")
	// form.Set("di", getDIHash(agent, profile))
	form.Set("nonce", nonce)
	form.Set("js", getJSHash(script, nonce, sessionID))
	form.Set("ai", getAIValue(script))
	// I'm guessing this means it's in prod
	form.Set("ii", "-1")
	form.Set("pi", getPIValue(script))
	form.Set("hk", "")
	form.Set("b", agent.BrowserName)
	form.Set("bv", agent.BrowserVersion)
	form.Set("bos", agent.OSName)
	form.Set("cb", "tdz_callback")
	// time in MS that it took to execute the script
	form.Set("et", fmt.Sprint(rand.Intn(30)+30))
	o := form.URLEncode()
	o = strings.ReplaceAll(o, "+", "%20")
	// xor the output with the nonce, return the base64 encoded result
	// if i > len(nonce), then we need to start at 0
	xor := ""
	for i := 0; i < len(o); i++ {
		xor += string(o[i] ^ nonce[i%len(nonce)])
	}
	a := base64.StdEncoding.EncodeToString([]byte(xor))
	return a
}

// this gets the dynamic pi value from the script
// TODO: fetch this dynamically, I cba to do that rn
func getPIValue(script string) string {
	return "99998"
}

// this gets the dynamic ai value from the script
// TODO: fetch this dynamically, I cba to do that rn
func getAIValue(script string) string {
	return "2212"
}

// this is the `di` hash
// it's just a sha-1 hash of a bunch of static browser functions outputs
/*
[
    " - navigator.vendorSub",
    "0 - navigator.maxTouchPoints",
    "1040 - screen.availHeight",
    "1080 - screen.height",
    "176f3ec286103cca16d5d7d410c073a0364ce9e4 - pluginsHash",
    "1920 - screen.availWidth",
    "1920 - screen.width",
    "20030107 - navigator.productSub",
    "24 - screen.colorDepth",
    "24 - screen.pixelDepth",
    "240 - new Date().getTimezoneOffset()",
    "32 - navigator.hardwareConcurrency",
    "40c8492954b1e2073ca612456ecb7be84ed4e72a - hash of fonts",
    "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 OPR/119.0.0.0 - navigator.appVersion",
    "Gecko - navigator.product",
    "Google Inc. - navigator.vendor",
    "Mozilla - navigator.appCodeName",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 OPR/119.0.0.0 - navigator.userAgent",
    "Netscape - navigator.appName",
    "Win32 - navigator.platform",
    "en-US - navigator.language",
    "false - navigator.webdriver",
    "true - navigator.cookieEnabled",
    "true - navigator.onLine"
]
*/
func getDIHash(agent *UserAgent, profile *Profile) string {
	out := []string{}
	// vendorSub
	out = append(out, "")
	// maxTouchPoints
	out = append(out, fmt.Sprintf("%d", profile.MaxTouchPoints))
	// availHeight
	out = append(out, fmt.Sprintf("%d", profile.ScreenAnvilHeight))
	// height
	out = append(out, fmt.Sprintf("%d", profile.ScreenHeight))
	// pluginsHash
	// hardcode pluginHash for now
	out = append(out, "176f3ec286103cca16d5d7d410c073a0364ce9e4")
	// availWidth
	out = append(out, fmt.Sprintf("%d", profile.ScreenAnvilWidth))
	// width
	out = append(out, fmt.Sprintf("%d", profile.ScreenWidth))
	// productSub
	// hardcode productSub for now
	out = append(out, "20030107")
	// colorDepth
	out = append(out, fmt.Sprintf("%d", profile.ColorDepth))
	// pixelDepth, this is the same as colorDepth
	out = append(out, fmt.Sprintf("%d", profile.ColorDepth))
	// timeZoneOffset
	out = append(out, fmt.Sprintf("%d", profile.TimezoneOffset))
	// hardwareConcurrency
	out = append(out, fmt.Sprintf("%d", profile.HardwareConcurrency))
	// hash of fonts
	out = append(out, "40c8492954b1e2073ca612456ecb7be84ed4e72a")
	// appVersion
	out = append(out, fmt.Sprintf("%s", agent.UA[8:]))
	// product
	out = append(out, "Gecko")
	// vendor
	out = append(out, "Google Inc.")
	// appCodeName
	out = append(out, "Mozilla")
	// userAgent
	out = append(out, agent.UA)
	// appName
	out = append(out, "Netscape")
	// platform
	out = append(out, "Win32")
	// language
	out = append(out, "en-US")
	// webdriver
	out = append(out, "false")
	// cookieEnabled
	out = append(out, "true")
	// onLine
	out = append(out, "true")
	// output is .join("") and sha-1 hashed
	hash := sha1.New()
	hash.Write([]byte(strings.Join(out, "")))
	return hex.EncodeToString(hash.Sum(nil))
}

// this is the FP hash
// this is supposed to be a hash of the HTML on the page
func getFPHash(htmlElements []string, nonce string) string {
	nonce += "000000000000000000000000"
	nonce = nonce[:24]
	hashes := []string{}
	h := sha1.New()
	for _, element := range htmlElements {
		xor := ""
		for i := 0; i < len(element); i++ {
			xor += string(element[i] ^ nonce[i%len(nonce)])
		}
		h.Reset()
		h.Write([]byte(xor))
		hashes = append(hashes, "0x"+hex.EncodeToString(h.Sum(nil)))
	}
	// sort the hashes
	sort.Strings(hashes)
	hash := sha1.New()
	hash.Write([]byte(strings.Join(hashes, "")))
	return hex.EncodeToString(hash.Sum(nil))
}

// all of this is the JS hash
func furtherFormatCode(td_Nj string, td_Rb bool) string {
	if td_Rb {
		td_Nj = regexp.MustCompile(` {2,}`).ReplaceAllString(td_Nj, " ")
		td_Nj = regexp.MustCompile(`[\n\r]*`).ReplaceAllString(td_Nj, "")
		td_Nj = strings.ReplaceAll(td_Nj, "\\\"", "")
		td_Nj = strings.ReplaceAll(td_Nj, "'", "")
		td_Nj = strings.ReplaceAll(td_Nj, `"`, "")
		td_Nj = strings.ReplaceAll(td_Nj, " ", "")
	} else {
		td_Nj = regexp.MustCompile(` {2,}`).ReplaceAllString(td_Nj, " ")
		td_Nj = regexp.MustCompile(`[\n\r\t;]*`).ReplaceAllString(td_Nj, "")
		td_Nj = strings.ReplaceAll(td_Nj, "'", "\\\"")
		td_Nj = strings.ReplaceAll(td_Nj, " ", "")
	}
	return td_Nj
}

type funcToGrab struct {
	Assigned bool
	Name     string
}

func getFuncHasher(funcs map[string]string) []*funcToGrab {
	funcsToGrab := []*funcToGrab{}
	for _, function := range funcs {
		// find the function we need
		if !(strings.Contains(function, `.join("");`) && strings.Contains(function, `.toString())`) && strings.Contains(function, `,true);`)) {
			continue
		}
		parts := strings.Split(strings.Split(strings.Split(function, "[")[1], "]")[0], ",")
		for _, p := range parts {
			p = strings.Replace(p, ".toString())", "", 1)
			p = strings.Split(p, "(")[1]
			funcsToGrab = append(funcsToGrab, &funcToGrab{
				Assigned: strings.Contains(p, "."),
				Name:     p,
			})
		}
		return funcsToGrab
	}
	return funcsToGrab
}

func getFuncBodies(code string) map[string]string {
	funcs := make(map[string]string)
	insideFunc := false
	conditionalDepth := 0
	funcCode := strings.Builder{}
	funcName := ""
	for i, f := range code {
		// fmt.Println(string(f))
		if f == 'n' {
			if code[i-7:i+1] == "function" {
				insideFunc = true
				// get the name of the function
				if code[i+1] == '(' {
					funcName = code[i-19 : i-8]
				} else {
					for j := i + 2; j < len(code); j++ {
						if code[j] != '(' {
							funcName += string(code[j])
							continue
						}

						break
					}
				}
				funcCode.WriteString(code[i-7 : i+1])
				conditionalDepth = 0
				continue
			}
		}
		if !insideFunc {
			continue
		}
		if insideFunc {
			funcCode.WriteRune(f)
		}
		if f == '{' {
			// make sure we're not hitting a string
			if code[i-1] == '"' || code[i-1] == '\'' {
				continue
			}
			conditionalDepth++
			continue
		}
		if f == '}' {
			// make sure we're not hitting a string
			if code[i-1] == '"' || code[i-1] == '\'' {
				continue
			}
			conditionalDepth--
			// function is done being visited
			if conditionalDepth == 0 {
				funcs[funcName] = funcCode.String()
				funcCode.Reset()
				insideFunc = false
				funcName = ""
				continue
			}
			continue
		}
	}
	return funcs
}

func getJSHash(script, nonce, sessionID string) string {
	// get the code of the functions
	/*
		how it works:
			1.) we check if the current byte is `n`
			2.) if it is, we check if the last 7 bytes are `functio`
			3.) if it is, we check if the i+2th byte is a `(`
			    - if it is, we know it's an assigned function
			4.) next to get the function body, we mark that we are inside a desirable function, so we start adding bytes, from i-7 to the ending braket
			    - if we hit a `}`, we need to know if we are inside a condition or not, if we aren't, then we know it's the end of the function
	*/
	funcs := getFuncBodies(script)
	// get the names of the functions we need to grab
	funcsToGet := getFuncHasher(funcs)
	// for i, f := range funcsToGet {
	// 	fmt.Printf("[%v] %v %s\n", i, f.Assigned, f.Name)
	// }
	// save all func strings to a single string in order
	code := ""
	for _, f := range funcsToGet {
		code += funcs[f.Name]
	}
	// now we further format the code
	code = furtherFormatCode(code, true)
	// sha1 hash
	hash := sha1.New()
	hash.Write([]byte(code))
	hashString := hex.EncodeToString(hash.Sum(nil))

	// get the second hash used for xor
	hash.Reset()
	// comment + body + nonce + 1 + 0
	hash.Write([]byte("#COMMENTBODY" + nonce + "10"))
	hash2String := hex.EncodeToString(hash.Sum(nil))
	// xor the two hashes
	xor := ""
	for i := 0; i < len(hashString); i++ {
		xor += string(hashString[i] ^ hash2String[i])
	}
	hash.Reset()
	// return a sha1 hash of the xor
	hash.Write([]byte(xor))
	out := "0x" + hex.EncodeToString(hash.Sum(nil))
	return out
}
