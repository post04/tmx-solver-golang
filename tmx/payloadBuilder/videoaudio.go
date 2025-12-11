package payloadbuilder

import (
	"fmt"
	"net/url"

	"github.com/obfio/tmx-solver-golang/mongo"
)

func videoAudio(p *mongo.Print) string {
	videoAudioPayload := &url.Values{}
	videoAudioPayload.Set("medh", fmt.Sprintf("(1,1,1,%s)", sha256Hex(p.General.MediaDevicesStr)))
	a := "&" + videoAudioPayload.Encode()
	return a
}
