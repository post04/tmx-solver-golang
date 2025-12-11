package geturls

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func getFuncBodies(code string) map[string]string {
	funcs := make(map[string]string)
	insideFunc := false
	conditionalDepth := 0
	funcCode := strings.Builder{}
	funcName := ""
	for i, f := range code {
		// fmt.Println(string(f))
		if f == 'n' && !insideFunc {
			if i-7 < 0 {
				continue
			}
			if code[i-7:i+1] == "function" {
				insideFunc = true
				// get the name of the function
				if code[i+1] == '(' {
					if i-19 < 0 {
						insideFunc = false
						continue
					}
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
	// add final func if it exists
	if funcName != "" {
		funcs[funcName] = funcCode.String()
	}
	return funcs
}

var (
	encodedStrRegex = regexp.MustCompile(`var td_[0-9A-z]{2}= new td_[0-9A-z]{2}\.td_[0-9A-z]{2}\("[A-z0-9]+"\)`)

	allDecodedStringsUsesRegex = regexp.MustCompile(`td_[0-9A-z]{2}=\(td_[0-9A-z]{2}\)\?td_[0-9A-z]{2}\.td_[0-9A-z]{1,2}\([0-9]+,[0-9]+\)`)

	// check.js URL regex
	checkJSURLRegex = regexp.MustCompile(`td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}\(td_[A-z0-9]{1,2},td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}`)

	// clear.png 1 URL regex
	clearPNG1URLRegex = regexp.MustCompile(`td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}\(td_[A-z0-9]{1,2},td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}\);`)

	// clear.png 2 URL regex
	clearPNG2URLRegex = regexp.MustCompile(`\+td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}\+`)

	// HPURL regex
	// .setAttribute("src", td_3M);
	HPURLRegex = regexp.MustCompile(`null;td_[A-z0-9]{1,2}\.setAttribute\(.+,td_[A-z0-9]{1,2}\);`)

	// blankClear.URL regex
	blankClearURLRegex = regexp.MustCompile(`td_[A-z0-9]{1,2}\.open\(.+,td_[A-z0-9]{1,2},true\);`)

	// blankClear.Nonce regex
	blankClearNonceRegex = regexp.MustCompile(`\+td_[A-z0-9]{1,2}\+td_[A-z0-9]{1,2};`)

	// ls_fp regexs
	lsFpIdentifierRegex     = regexp.MustCompile(`if\(td_[A-z0-9]{1,2}!==null\){td_[A-z0-9]{1,2}\+=`)
	lsFpFuncIdentifierRegex = regexp.MustCompile(`if\(typeoftd_[A-z0-9]{1,2}!==\[\]\[\[\]\]\+""\){td_[A-z0-9]{1,2}\(\);}`)
	lsFpRegex               = regexp.MustCompile(`td_[A-z0-9]{1,2}\.load_iframe\(td_[A-z0-9]{1,2},document\);`)

	// LSA regexes
	lsaURLRegex   = regexp.MustCompile(`vartd_[A-z0-9]{1,2}=td_[A-z0-9]{1,2}\+td_[A-z0-9]{1,2};`)
	lsaNonceRegex = regexp.MustCompile(`[A-z0-9]{2}\.split\("_"\)\[0\];`)

	// top_fp.html regexes
	topFpHTMLRegex = regexp.MustCompile(`functiontd_[A-z0-9]{1,2}\(\)\{td_[A-z0-9]{1,2}\.load_iframe\(td_[A-z0-9]{1,2},document\);\}`)

	// h64 regexes
	h64IdentifierRegex = regexp.MustCompile(`else\{if\(typeoftd_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}===`)
	h64Regex           = regexp.MustCompile(`\}td_[A-z0-9]{1,2}=td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2};`)

	// long subdomain URL regexes
	longSubdomainIdentifierRegex = regexp.MustCompile(`td_[A-z0-9]{1,2}=td_[A-z0-9]{1,2}\(\);td_[A-z0-9]{1,2}\(\);`)

	// audio/video URL regexes
	videoAudioURLRegex = regexp.MustCompile(`navigator\.sendBeacon\(td_[A-z0-9]{1,2},td_[A-z0-9]{1,2}\);`)

	// sid URL regexes
	sidURLRegex   = regexp.MustCompile(`td_[A-z0-9]{1,2}\.td_[A-z0-9]{1,2}\(td_[A-z0-9]{1,2},td_[A-z0-9]{1,2}\);td_[A-z0-9]{1,2}\(td_[A-z0-9]{1,2},document\);\}`)
	sidURLMatcher = regexp.MustCompile(`vartd_[A-z0-9]{1,2}=td_[A-z0-9]{1,2}\+`)

	// sid nonce regex
	sidNonceRegex = regexp.MustCompile(`\{vartd_[A-z0-9]{1,2}=td_[A-z0-9]{1,2}\?td_[A-z0-9]{1,2}\:td_[A-z0-9]{1,2};`)

	// ARF URL regex
	arfURLRegex = regexp.MustCompile(`td_[A-z0-9]{1,2}\(td_[A-z0-9]{1,2},td_[A-z0-9]{1,2},td_[A-z0-9]{1,2}\.join\(""\)\);\}`)
)

func GetClosestStringDecoder(decoderPositions [][]int, firstURLPos []int) int {
	smallestDiff := math.MaxInt
	closestDecoderPos := 0
	for i := 0; i < len(decoderPositions); i++ {
		diff := decoderPositions[i][1] - firstURLPos[0]
		if !(diff < 0) {
			continue
		}
		if diff < smallestDiff {
			smallestDiff = -diff
			closestDecoderPos = i
		}
	}
	return closestDecoderPos
}

func getEncodingPositions(script string) ([][]int, error) {
	positions := encodedStrRegex.FindAllStringIndex(script, -1)
	if len(positions) == 0 {
		return nil, errors.New("encodingPositions is empty")
	}
	return positions, nil
}

func getDecodedStringsUses(script string) ([][]int, error) {
	positions := allDecodedStringsUsesRegex.FindAllStringIndex(script, -1)
	if len(positions) == 0 {
		return nil, errors.New("decodedStringsUses is empty")
	}
	return positions, nil
}

func getPossibleEncodingStrs(encodingPositions [][]int, script string) ([]string, error) {
	possibleEncodingStrs := make([]string, len(encodingPositions))
	for i := 0; i < len(encodingPositions); i++ {
		possibleEncodingStr := script[encodingPositions[i][0]:encodingPositions[i][1]]
		possibleEncodingStr = strings.Split(possibleEncodingStr, "(")[1]
		possibleEncodingStr = strings.Split(possibleEncodingStr, ")")[0]
		possibleEncodingStr = possibleEncodingStr[1 : len(possibleEncodingStr)-1]
		possibleEncodingStrs[i] = possibleEncodingStr
	}
	return possibleEncodingStrs, nil
}

func extractDecodedNumbers(variable string, decodedStringsUses [][]int, decodedStringsUsesStrs []string, encodingPositions [][]int, possibleEncodingStrs []string, script string) ([2]int, string, error) {
	// Find the position of the variable in decodedStringsUses
	variablePos := []int{}
	for i := 0; i < len(decodedStringsUsesStrs); i++ {
		if strings.Contains(decodedStringsUsesStrs[i], variable) {
			variablePos = decodedStringsUses[i]
			break
		}
	}
	if len(variablePos) == 0 {
		return [2]int{}, "", fmt.Errorf("variable %s position not found", variable)
	}

	// Get the closest decoder
	closestDecoderPos := GetClosestStringDecoder(encodingPositions, variablePos)
	decoderStr := possibleEncodingStrs[closestDecoderPos]

	// Extract the numbers from the script
	variableSection := script[variablePos[0]:variablePos[1]]
	numStrs := strings.Split(strings.Split(strings.Split(strings.Split(variableSection, "?")[1], "(")[1], ")")[0], ",")

	if len(numStrs) != 2 {
		return [2]int{}, "", fmt.Errorf("expected 2 numbers for variable %s, got %d", variable, len(numStrs))
	}

	nums := [2]int{}
	for i := 0; i < 2; i++ {
		num, err := strconv.Atoi(numStrs[i])
		if err != nil {
			return [2]int{}, "", fmt.Errorf("failed to convert string to int for variable %s: %w", variable, err)
		}
		nums[i] = num
	}

	return nums, decoderStr, nil
}
