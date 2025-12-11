package geturls

import (
	"fmt"
	"strings"

	"github.com/obfio/tmx-solver-golang/tmx/decodeObfuscation"
)

type SidFPDynamic struct {
	SID *sid
}

func GetSidFPDynamic(script string) (*SidFPDynamic, error) {
	tag := &SidFPDynamic{}

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
			sidNonceVariable = strings.Split(sidNonceVariable, ":")[0]
			sidNonceVariable = strings.Split(sidNonceVariable, "?")[1]

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

	return tag, nil
}
