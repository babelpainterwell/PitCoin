package hashutil

import "crypto/sha256"

func DoubleSha256(data []byte) [32]byte {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}