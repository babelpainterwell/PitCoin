package hashutil

import (
	"crypto/sha256"
)

func DoubleSha256(data []byte) [32]byte {
	// sha256.Sum256 expects a byte slice and returns an array of 32 bytes

	if data == nil {
		return [32]byte{}
	}

	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}


// func main() {
// 	// testing 

// 	data1 := []byte("5")
// 	data2 := []byte("6")
// 	hash_1 := DoubleSha256(data1)
// 	hash_2 := DoubleSha256(data2)

// 	F := doubleSha256Concat(hash_1, hash_2)

// 	C_F := DoubleSha256(F[:])
// 	C_F_F := doubleSha256Concat(F, F)

// 	fmt.Println(C_F)
// 	fmt.Println(C_F_F)

// }