package geturls

import (
	"fmt"
	"strings"

	"github.com/obfio/tmx-solver-golang/tmx/decodeObfuscation"
)

type ARFDynamic struct {
	ARFURL   string
	ARFNonce string
}

func GetARFDynamic(script string) (*ARFDynamic, error) {
	tag := &ARFDynamic{}

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

	for _, funcBody := range funcs {
		if arfURLRegex.MatchString(funcBody) {
			arfURLVariable := arfURLRegex.FindString(funcBody)
			arfURLVariable = strings.Split(arfURLVariable, "(")[1]
			arfURLVariable = strings.Split(arfURLVariable, ",")[0]

			nums, decoderStr, err := extractDecodedNumbers(arfURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract ARF URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.ARFURL = decoder.Decode(nums[0], nums[1])
			// arf nonce now
			arfNonceVariable := arfURLRegex.FindString(funcBody)
			arfNonceVariable = strings.Split(arfNonceVariable, "(")[1]
			arfNonceVariable = strings.Split(arfNonceVariable, ",")[1]
			nums, decoderStr, err = extractDecodedNumbers(arfNonceVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract ARF URL numbers: %w", err)
			}
			decoder = decodeObfuscation.CreateDecoder(decoderStr)
			tag.ARFNonce = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	return tag, nil
}
