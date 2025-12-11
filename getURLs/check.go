package geturls

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/obfio/tmx-solver-golang/tmx/decodeObfuscation"
	"github.com/t14raptor/go-fast/ast"
	"github.com/t14raptor/go-fast/generator"
	"github.com/t14raptor/go-fast/parser"
)

type lsa struct {
	URL   string
	Nonce string
}

type blankClear struct {
	URL   string
	Nonce string
}

type sid struct {
	URL   string
	Nonce string
}

type CheckJSDynamic struct {
	// not required, don't error on this
	HPURL string

	BlankClear *blankClear

	LSFPURL string

	LSA *lsa

	SidFPURL string

	ESJSURL string

	TOPFPHTMLURL string

	JAJBURL string

	H64URL string

	LongSubdomainURL string

	VIDEOAUDIOURL string

	WGLURL string

	IPURL string

	SID *sid
	// probably more to do after this
}

type FunctionDeclarationVisitor struct {
	ast.NoopVisitor
	functions map[string]string
}

func (v *FunctionDeclarationVisitor) VisitFunctionDeclaration(node *ast.FunctionDeclaration) {
	v.functions[node.Function.Name.Name] = generator.Generate(node)
}

func GetCheckJSDynamic(script string) (*CheckJSDynamic, error) {

	tag := &CheckJSDynamic{}

	encodingPositions, err := getEncodingPositions(script)
	if err != nil {
		return nil, err
	}
	possibleEncodingStrs, err := getPossibleEncodingStrs(encodingPositions, script)
	if err != nil {
		return nil, err
	}

	// get all possible definitons of decoded strings
	decodedStringsUses, err := getDecodedStringsUses(script)
	if err != nil {
		return nil, err
	}
	decodedStringsUsesStrs := make([]string, len(decodedStringsUses))
	for i := 0; i < len(decodedStringsUses); i++ {
		decodedStringsUsesStrs[i] = script[decodedStringsUses[i][0]:decodedStringsUses[i][1]]
	}

	// get all the functions defined in the script
	funcs := make(map[string]string)
	// use go-fast to parse the script
	parser, err := parser.ParseFile(script)
	if err != nil {
		return nil, err
	}
	visitor := &FunctionDeclarationVisitor{
		functions: make(map[string]string),
	}
	visitor.V = visitor
	parser.VisitWith(visitor)
	funcs = visitor.functions
	for name, funcBody := range funcs {
		funcBody = strings.ReplaceAll(funcBody, "\n", "")
		funcBody = strings.ReplaceAll(funcBody, "\r", "")
		funcBody = strings.ReplaceAll(funcBody, "\t", "")
		funcBody = strings.ReplaceAll(funcBody, " ", "")
		funcs[name] = funcBody
	}

	// a, err := json.MarshalIndent(funcs, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// os.WriteFile("funcs.json", a, 0644)

	// get HPURL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, ".onload=null;") {

			HPURLUsage := HPURLRegex.FindString(funcBody)
			if HPURLUsage == "" {
				fmt.Println("FAILED ON HP")
				os.WriteFile("FAILED.js", []byte(script), 0666)
				break
			}
			parts := strings.Split(HPURLUsage, ",")
			HPURLVariable := parts[len(parts)-1]
			HPURLVariable = strings.Split(HPURLVariable, ")")[0]
			nums, decoderStr, err := extractDecodedNumbers(HPURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				fmt.Println("FAILED ON HP (1)")
				// HPURL is optional, don't error on this
				break
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.HPURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	blank := &blankClear{}
	// get BlankClear URL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, "replace(/[\\r\\n]/g,\"\");") {
			blankClearURLUsage := blankClearURLRegex.FindString(funcBody)
			if blankClearURLUsage == "" {
				return nil, errors.New("blankClearURLUsage is empty")
			}
			blankClearURLVariable := strings.Split(blankClearURLUsage, ",true")[0]
			blankClearURLVariable = blankClearURLVariable[len(blankClearURLVariable)-5:]

			nums, decoderStr, err := extractDecodedNumbers(blankClearURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract BlankClear URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			blank.URL = decoder.Decode(nums[0], nums[1])
			break
		}
	}
	// if blank.URL == "" {
	// 	os.WriteFile("FAILED.js", []byte(script), 0666)
	// 	bb, _ := json.MarshalIndent(funcs, "", "	")
	// 	os.WriteFile("funcs.json", bb, 0666)
	// }

	// get BlankClear nonce
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, "replace(/[\\r\\n]/g,\"\");") {
			blankClearNonceUsage := blankClearNonceRegex.FindString(funcBody)
			if blankClearNonceUsage == "" {
				return nil, errors.New("blankClearNonceUsage is empty")
			}
			blankClearNonceVariable := strings.Split(blankClearNonceUsage, "+")[1]

			nums, decoderStr, err := extractDecodedNumbers(blankClearNonceVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract BlankClear Nonce numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			blank.Nonce = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	tag.BlankClear = blank
	// aa, _ := json.MarshalIndent(tag, "", "	")
	// panic(string(aa))

	for _, funcBody := range funcs {
		if lsFpIdentifierRegex.MatchString(funcBody) {
			// get the name of the function that we need
			lsFpFuncName := lsFpFuncIdentifierRegex.FindString(funcBody)
			lsFpFuncName = strings.Split(lsFpFuncName, "{")[1]
			lsFpFuncName = strings.Split(lsFpFuncName, "()")[0]
			funcBody = funcs[lsFpFuncName]
			// if the funcBody is empty, try to get the funcBody from another function.
			for _, fb := range funcs {
				if strings.Contains(fb, "function"+lsFpFuncName) {
					funcBody = fb
					break
				}
			}
			if funcBody == "" {
				return nil, errors.New("lsFpFunc Body is empty")
			}
			lsFpURLUsage := lsFpRegex.FindString(funcBody)
			if lsFpURLUsage == "" {
				return nil, errors.New("lsFpURLUsage is empty")
			}
			lsFpURLVariable := strings.Split(lsFpURLUsage, ",")[0]
			lsFpURLVariable = strings.Split(lsFpURLVariable, "(")[1]

			nums, decoderStr, err := extractDecodedNumbers(lsFpURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract lsFp URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.LSFPURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	lsa := &lsa{}
	// get LSA URL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, "if(window.localStorage){") {
			lsaURLUsage := lsaURLRegex.FindString(funcBody)
			if lsaURLUsage == "" {
				return nil, errors.New("lsaURLUsage is empty")
			}
			lsaURLVariable := strings.Split(lsaURLUsage, "=")[1]
			lsaURLVariable = strings.Split(lsaURLVariable, "+")[0]

			nums, decoderStr, err := extractDecodedNumbers(lsaURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract LSA URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			lsa.URL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// get LSA nonce
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, "if(window.localStorage){") {
			lsaNonceUsage := lsaNonceRegex.FindString(funcBody)
			if lsaNonceUsage == "" {
				return nil, errors.New("lsaNonceUsage is empty")
			}
			lsaNonceVariable := strings.Split(lsaNonceUsage, ".")[0]

			nums, decoderStr, err := extractDecodedNumbers(lsaNonceVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract LSA Nonce numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)

			lsa.Nonce = decoder.Decode(nums[0], nums[1])
			break
		}
	}
	tag.LSA = lsa

	// sid_fp URL
	for _, funcBody := range funcs {
		if lsFpIdentifierRegex.MatchString(funcBody) {
			// get the name of the function that we need
			sidFpFuncNames := lsFpFuncIdentifierRegex.FindAllString(funcBody, -1)
			sidFpFuncName := strings.Split(sidFpFuncNames[1], "{")[1]
			sidFpFuncName = strings.Split(sidFpFuncName, "()")[0]
			funcBody = funcs[sidFpFuncName]
			// if the funcBody is empty, try to get the funcBody from another function.
			for _, fb := range funcs {
				if strings.Contains(fb, "function"+sidFpFuncName) {
					funcBody = fb
					break
				}
			}
			if funcBody == "" {
				tag.SidFPURL = "FAILED"
				break
			}
			sidFpURLUsage := lsFpRegex.FindString(funcBody)
			if sidFpURLUsage == "" {
				bbb, _ := json.MarshalIndent(funcs, "", "	")
				os.WriteFile("funcs.json", bbb, 0666)
				return nil, errors.New("sidFpURLUsage is empty")
			}
			sidFpURLVariable := strings.Split(sidFpURLUsage, ",")[0]
			sidFpURLVariable = strings.Split(sidFpURLVariable, "(")[1]

			nums, decoderStr, err := extractDecodedNumbers(sidFpURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract sidFp URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.SidFPURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// es.js URL
	for _, funcBody := range funcs {
		if !strings.Contains(funcBody, "if(window.localStorage){") && strings.Contains(funcBody, "window.localStorage.getItem") {
			esJsURLVariable := strings.Split(funcBody, "===")[0]
			esJsURLVariable = strings.Split(esJsURLVariable, "typeof")[1]

			nums, decoderStr, err := extractDecodedNumbers(esJsURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract esJs URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.ESJSURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// top_fp.html URL
	for _, funcBody := range funcs {
		if topFpHTMLRegex.MatchString(funcBody) {
			topFpURLVariable := strings.Split(funcBody, ",")[0]
			topFpURLVariable = strings.Split(topFpURLVariable, "(")[2]
			nums, decoderStr, err := extractDecodedNumbers(topFpURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract topFp URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.TOPFPHTMLURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// jajb URL
	tag.JAJBURL = tag.LSA.URL

	// h64 URL
	for _, funcBody := range funcs {
		if h64IdentifierRegex.MatchString(funcBody) {
			h64URLVariable := h64Regex.FindString(funcBody)
			h64URLVariable = strings.Split(h64URLVariable, ".")[1]
			h64URLVariable = strings.Split(h64URLVariable, ";")[0]

			nums, decoderStr, err := extractDecodedNumbers(h64URLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract h64 URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.H64URL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// long subdomain URL
	for _, funcBody := range funcs {
		if lsFpIdentifierRegex.MatchString(funcBody) {
			// get the name of the function that we need
			longSubdomainFuncNames := longSubdomainIdentifierRegex.FindAllString(funcBody, -1)

			longSubdomainFuncName := strings.Split(longSubdomainFuncNames[1], ";")[1]
			longSubdomainFuncName = strings.Split(longSubdomainFuncName, "()")[0]

			funcBody = funcs[longSubdomainFuncName]
			longSubdomainURLVariable := strings.Split(funcBody, ",")[0]
			longSubdomainURLVariable = strings.Split(longSubdomainURLVariable, "(")[2]

			nums, decoderStr, err := extractDecodedNumbers(longSubdomainURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract longSubdomain URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.LongSubdomainURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// video/audio URL
	for _, funcBody := range funcs {
		if videoAudioURLRegex.MatchString(funcBody) {
			videoAudioURLUsage := videoAudioURLRegex.FindString(funcBody)
			videoAudioURLVariable := strings.Split(videoAudioURLUsage, "(")[1]
			videoAudioURLVariable = strings.Split(videoAudioURLVariable, ",")[0]

			nums, decoderStr, err := extractDecodedNumbers(videoAudioURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract video/audio URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.VIDEOAUDIOURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// WGL URL
	tag.WGLURL = tag.JAJBURL

	sid := &sid{}

	// get SID URL
	for _, funcBody := range funcs {
		if sidURLRegex.MatchString(funcBody) && funcBody[16] != '{' {
			sidURLVariable := sidURLMatcher.FindString(funcBody)
			sidURLVariable = strings.Split(sidURLVariable, "=")[1]
			sidURLVariable = strings.Split(sidURLVariable, "+")[0]

			nums, decoderStr, err := extractDecodedNumbers(sidURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract SID URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			sid.URL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// get SID nonce
	for _, funcBody := range funcs {
		if sidNonceRegex.MatchString(funcBody) {
			sidNonceVariable := sidNonceRegex.FindString(funcBody)
			sidNonceVariable = strings.Split(sidNonceVariable, ":")[1]
			sidNonceVariable = strings.Split(sidNonceVariable, ";")[0]

			nums, decoderStr, err := extractDecodedNumbers(sidNonceVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract SID nonce numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			sid.Nonce = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	tag.SID = sid

	tag.IPURL = tag.JAJBURL

	// a, _ := json.MarshalIndent(funcs, "", "	")
	// os.WriteFile("funcs.json", a, 0666)
	// fmt.Println(string(a))
	// sanity check, make sure all the URLs and nonces are not empty. If they are, error saying what was missing
	if tag.BlankClear.URL == "" {
		return nil, errors.New("BlankClear URL is empty")
	}
	if tag.BlankClear.Nonce == "" {
		return nil, errors.New("BlankClear Nonce is empty")
	}
	if tag.LSFPURL == "" {
		return nil, errors.New("LSFPURL is empty")
	}
	if tag.LSA.URL == "" {
		return nil, errors.New("LSA URL is empty")
	}
	if tag.LSA.Nonce == "" {
		return nil, errors.New("LSA Nonce is empty")
	}
	if tag.SidFPURL == "" {
		return nil, errors.New("SidFPURL is empty")
	}
	if tag.ESJSURL == "" {
		return nil, errors.New("ESJSURL is empty")
	}
	if tag.TOPFPHTMLURL == "" {
		return nil, errors.New("TOPFPHTMLURL is empty")
	}
	if tag.JAJBURL == "" {
		return nil, errors.New("JAJBURL is empty")
	}
	if tag.H64URL == "" {
		return nil, errors.New("H64URL is empty")
	}
	if tag.LongSubdomainURL == "" {
		return nil, errors.New("LongSubdomainURL is empty")
	}
	if tag.VIDEOAUDIOURL == "" {
		return nil, errors.New("VIDEOAUDIOURL is empty")
	}
	if tag.WGLURL == "" {
		return nil, errors.New("WGLURL is empty")
	}
	if tag.IPURL == "" {
		return nil, errors.New("IPURL is empty")
	}
	if tag.SID.URL == "" {
		return nil, errors.New("SID URL is empty")
	}
	if tag.SID.Nonce == "" {
		os.WriteFile("FAILED.js", []byte(script), 0666)
		return nil, errors.New("SID Nonce is empty")
	}

	return tag, nil
}
