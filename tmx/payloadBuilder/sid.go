package payloadbuilder

import (
	"fmt"
	"strings"
	"time"

	"github.com/267H/orderedform"
)

func sid(r *Request) string {

	timeSeconds := fmt.Sprint(time.Now().UnixMilli() / 1000)
	rnd := "tdr_" + randString(16)
	form := orderedform.NewForm(6)
	form.Set("sid_rnd", rnd)
	form.Set("sid_date", timeSeconds)
	form.Set("sid_type", "web:ecdsa")
	enc := generateEncryptedFingerprint(timeSeconds, rnd, r.Nonce)
	form.Set("sid_key", enc.PubKeyEncoded)
	form.Set("sid_sig", enc.Output)
	form.Set("sifr", "0")
	a := strings.ReplaceAll(form.URLEncode(), "%3A", ":")
	return a

}
