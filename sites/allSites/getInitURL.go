package allsites

import (
	"math"
	"math/rand"
	"strconv"
)

func randomShit() int {
	a := rand.Float64()
	a = 2742745743359.0 * a
	a = math.Floor(a)
	return int(a)
}

func toString(e int) string {
	return strconv.FormatInt(int64(e+78364164096), 36)
}

func getRandomStr() string {
	return toString(randomShit())
}

func getInitURL(siteID, pageID, URL, sessionID string) string {
	c := randomShit()
	b := getRandomStr() + toString(c)
	// only set pageID if it's not empty
	URLOutput := URL + b + ".js?" + getRandomStr() + getRandomStr() + "=" + siteID + "&" + getRandomStr() + getRandomStr() + "=" + sessionID
	if pageID != "" {
		URLOutput += "&" + getRandomStr() + getRandomStr() + "=" + pageID
	}
	return URLOutput
}
