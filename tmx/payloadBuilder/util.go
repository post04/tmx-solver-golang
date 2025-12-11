package payloadbuilder

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"

	cryptoRand "crypto/rand"
)

type encFP struct {
	PubKeyEncoded string
	Output        string
}

func generateEncryptedFingerprint(timeSeconds, rnd, nonce string) *encFP {
	out := &encFP{}
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptoRand.Reader)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return nil
	}
	derBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		return nil
	}
	hexString := hex.EncodeToString(derBytes)
	out.PubKeyEncoded = hexString
	message := rnd + nonce + timeSeconds + "web:ecdsa"
	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(cryptoRand.Reader, privKey, hash[:])
	if err != nil {
		fmt.Println("Error signing message:", err)
		return nil
	}

	// Manually encode the signature to DER format.
	derSignature := encodeSignature(r, s)

	// Convert the DER-encoded signature to a hex string.
	hexSignature := hex.EncodeToString(derSignature)
	out.Output = hexSignature
	return out
}

func randString(length int) string {
	out := ""
	for i := 0; i < length; i++ {
		randInt := rand.Intn(62)
		if randInt < 10 {
			out += fmt.Sprint(randInt)
			continue
		}
		if randInt < 36 {
			out += fromCharCode(randInt + 55)
			continue
		}
		out += fromCharCode(randInt + 61)
	}
	return out
}

func sha256Hex(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func md5Hex(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// FromCharCode https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/FromCharCode
func fromCharCode(c int) string {
	return string(rune(c))
}

func encodeSignature(r, s *big.Int) []byte {
	// Encode each integer (r and s) individually.
	rEncoded := encodeInteger(r)
	sEncoded := encodeInteger(s)

	// Concatenate the two INTEGER encodings.
	sequenceContent := append(rEncoded, sEncoded...)

	// Wrap the concatenated content in a SEQUENCE.
	// 0x30 is the tag for a SEQUENCE.
	der := []byte{0x30, byte(len(sequenceContent))}
	der = append(der, sequenceContent...)
	return der
}

// encodeInteger encodes a big.Int into its minimal DER INTEGER form.
// It performs similar steps to the JavaScript code:
// 1. Obtains the minimal big-endian representation (which removes extra leading zeros).
// 2. If the first byte is >= 0x80, it prepends a 0x00 byte to indicate a positive integer.
// 3. Wraps the result with the INTEGER tag (0x02) and its length.
func encodeInteger(i *big.Int) []byte {
	// Get the minimal big-endian representation.
	b := i.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}

	// If the most significant bit is set, prepend a zero byte.
	if b[0] >= 0x80 {
		b = append([]byte{0x00}, b...)
	}

	// Wrap with the INTEGER tag (0x02) and the length.
	result := []byte{0x02, byte(len(b))}
	result = append(result, b...)
	return result
}
