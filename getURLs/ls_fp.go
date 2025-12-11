package geturls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/obfio/tmx-solver-golang/tmx/decodeObfuscation"
)

type lsb struct {
	URL   string
	Nonce string
}

type LsFpDynamic struct {
	LSB *lsb
}

func GetLsFpDynamic(script string) (*LsFpDynamic, error) {
	tag := &LsFpDynamic{}
	script = strings.ReplaceAll(script, `<html lang="en"><title>empty</title><body><script type="text/javascript">`, "")
	script = strings.ReplaceAll(script, `</script></body></html>`, "")
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
	funcs := getFuncBodies(script)
	for name, funcBody := range funcs {
		funcBody = strings.ReplaceAll(funcBody, "\n", "")
		funcBody = strings.ReplaceAll(funcBody, "\r", "")
		funcBody = strings.ReplaceAll(funcBody, "\t", "")
		funcBody = strings.ReplaceAll(funcBody, " ", "")
		funcs[name] = funcBody
	}

	lsa := &lsb{}
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
	tag.LSB = lsa
	return tag, nil
}
