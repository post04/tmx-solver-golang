package sites_mobile

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"hash/crc32"
)

// ALL AI GENERATED

// deinterleave separates the interleaved data into content and nonce
func deinterleave(data []byte, nonceSize int) (content []byte, nonce []byte, err error) {
	if len(data) < nonceSize {
		return nil, nil, fmt.Errorf("data too short: %d < %d", len(data), nonceSize)
	}

	contentBuffer := bytes.NewBuffer(nil)
	nonceBuffer := bytes.NewBuffer(nil)

	// Handle interleaving
	maxIndex := 2 * nonceSize
	if len(data) < maxIndex {
		maxIndex = len(data)
	}

	for i := 0; i < maxIndex; i += 2 {
		if i+1 < len(data) {
			contentBuffer.WriteByte(data[i])
			nonceBuffer.WriteByte(data[i+1])
		}
	}

	// Handle any remaining data
	for i := maxIndex; i < len(data); i++ {
		contentBuffer.WriteByte(data[i])
	}

	return contentBuffer.Bytes(), nonceBuffer.Bytes(), nil
}

// generateKey creates a 16-byte decryption key based on nonce, orgID, and sessionID
func generateKey(nonce []byte, orgID, sessionID string) []byte {
	key := make([]byte, 16)

	// Calculate the four CRC32 values as in the Java code
	orgIDBytes := []byte(orgID)
	sessionIDBytes := []byte(sessionID)

	// CRC32 of nonce + sessionID
	crc1 := crc32Value(append(nonce, sessionIDBytes...))

	// CRC32 of nonce + orgID
	crc2 := crc32Value(append(nonce, orgIDBytes...))

	// CRC32 of sessionID + nonce
	crc3 := crc32Value(append(sessionIDBytes, nonce...))

	// CRC32 of orgID + nonce
	crc4 := crc32Value(append(orgIDBytes, nonce...))

	// Combine the CRC32 values into a 16-byte key
	key[0] = byte(crc1 >> 24)
	key[1] = byte(crc1 >> 16)
	key[2] = byte(crc1 >> 8)
	key[3] = byte(crc1)

	key[4] = byte(crc2 >> 24)
	key[5] = byte(crc2 >> 16)
	key[6] = byte(crc2 >> 8)
	key[7] = byte(crc2)

	key[8] = byte(crc3 >> 24)
	key[9] = byte(crc3 >> 16)
	key[10] = byte(crc3 >> 8)
	key[11] = byte(crc3)

	key[12] = byte(crc4 >> 24)
	key[13] = byte(crc4 >> 16)
	key[14] = byte(crc4 >> 8)
	key[15] = byte(crc4)

	return key
}

// crc32Value calculates a CRC32 checksum
func crc32Value(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// decrypt decrypts the encrypted data using AES-CTR mode
func decrypt(encrypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a stream cipher using CTR mode
	stream := cipher.NewCTR(block, iv)

	// Decrypt the data
	decrypted := make([]byte, len(encrypted))
	stream.XORKeyStream(decrypted, encrypted)

	return decrypted, nil
}

// DecryptTMXPayload decrypts a base64-encoded TMX payload using the organization ID and session ID
func DecryptTMXPayload(encodedPayload, orgID, sessionID string) (string, error) {
	// Decode the base64 payload
	data, err := base64.StdEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	// Ensure the data is long enough
	if len(data) < 32 {
		return "", fmt.Errorf("data too short: %d < 32", len(data))
	}

	// Deinterleave the data to get the content and nonce (IV)
	content, nonce, err := deinterleave(data, 16)
	if err != nil {
		return "", err
	}

	// Generate the decryption key
	key := generateKey(nonce, orgID, sessionID)

	// Decrypt the content
	decrypted, err := decrypt(content, key, nonce)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
