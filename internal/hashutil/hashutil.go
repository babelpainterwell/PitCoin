package hashutil

import "crypto/sha256"

func DoubleSha256(data []byte) [32]byte {
	// sha256.Sum256 expects a byte slice and returns an array of 32 bytes
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}