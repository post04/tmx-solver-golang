package kleinanzeigenMobile

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

func makeAndroidID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func getSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func getGUID(androidID string) string {
	if len(androidID) == 0 {
		return ""
	}

	if len(androidID) == 32 {
		return androidID
	}

	md5Hash := md5.Sum([]byte(androidID))
	h := hex.EncodeToString(md5Hash[:])
	androidID += h[0 : 32-len(androidID)]

	return androidID
}
