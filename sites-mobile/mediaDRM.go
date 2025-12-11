package sites_mobile

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateDeviceUniqueId() string {
	const deviceUniqueIdSize = 16 // typically 16 bytes for a unique ID
	randomBytes := make([]byte, deviceUniqueIdSize)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(randomBytes)
}
