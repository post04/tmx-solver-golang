package geturls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/obfio/tmx-solver-golang/tmx/decodeObfuscation"
)

type TagDynamic struct {
	CheckJSURL   string
	ClearPNGURL  string
	ClearPNG2URL string
}

func GetTagDynamic(script string) (*TagDynamic, error) {
	tag := &TagDynamic{}

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
	// a, err := json.MarshalIndent(funcs, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// os.WriteFile("funcs.json", a, 0644)

	// get check.js URL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, ".injected=true;") {
			checkJSURLUsage := checkJSURLRegex.FindString(funcBody)
			if checkJSURLUsage == "" {
				return nil, errors.New("checkJSURLUsage is empty")
			}
			checkJSURLVariable := strings.Split(checkJSURLUsage, ".")[2]

			nums, decoderStr, err := extractDecodedNumbers(checkJSURLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract checkJS URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.CheckJSURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// get clear.png URL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, ".style.visibility") {
			clearPNG1URLUsage := clearPNG1URLRegex.FindString(funcBody)
			if clearPNG1URLUsage == "" {
				return nil, errors.New("clearPNG1URLUsage is empty")
			}
			clearPNG1URLVariable := strings.Split(clearPNG1URLUsage, ".")[2]
			clearPNG1URLVariable = clearPNG1URLVariable[:len(clearPNG1URLVariable)-2]

			nums, decoderStr, err := extractDecodedNumbers(clearPNG1URLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract clear PNG URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.ClearPNGURL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	// get clear.png 2 URL
	for _, funcBody := range funcs {
		if strings.Contains(funcBody, ".style.visibility") {
			clearPNG2URLUsage := clearPNG2URLRegex.FindString(funcBody)
			if clearPNG2URLUsage == "" {
				return nil, errors.New("clearPNG2URLUsage is empty")
			}
			clearPNG2URLVariable := strings.Split(clearPNG2URLUsage, ".")[1]
			clearPNG2URLVariable = clearPNG2URLVariable[:len(clearPNG2URLVariable)-1]

			nums, decoderStr, err := extractDecodedNumbers(clearPNG2URLVariable, decodedStringsUses, decodedStringsUsesStrs, encodingPositions, possibleEncodingStrs, script)
			if err != nil {
				return nil, fmt.Errorf("failed to extract clear PNG 2 URL numbers: %w", err)
			}

			decoder := decodeObfuscation.CreateDecoder(decoderStr)
			tag.ClearPNG2URL = decoder.Decode(nums[0], nums[1])
			break
		}
	}

	return tag, nil
}
