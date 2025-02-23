package main

// a minimal demo of proof of work consensus

import (
	"bytes"
	"fmt"
	"time"

	"github.com/babelpainterwell/shitcoin/internal/block"
)

func main() {

	fmt.Println("--- Minimal PoW Demo ---")

	var demoHeader block.BlockHeader

	// set the block header fields
	demoHeader.Version = 1
	demoHeader.Timestamp = uint32(time.Now().Unix())

	
	// define a difficulty target as 16 bits == 2 bytes, 16 zeros in a 256-bit hash
	// the number of leading zeros in the 32-byte hash, but might be too coarse for real mining
	difficulty := 2 
	targetPrefix := make([]byte, difficulty)

	tries := 0 
	startTime := time.Now()

	// mining 
	for {
		tries ++ 

		// increment the nonce field
		// the hash field is recomputed for each new nonce 
		currHeaderHash := demoHeader.BlockHash() 

		if bytes.HasPrefix(currHeaderHash[:], targetPrefix) {
			fmt.Println("--- Found a valid nonce ---")
			fmt.Println("Nonce: ", demoHeader.Nonce)
			fmt.Println("Hash (little endian): ", currHeaderHash)
			fmt.Println("Tries: ", tries)
			break
		}

		// increment the nonce field
		// we don't update the timestamp field, even in real Bitcoin mining
		demoHeader.Nonce++

		if tries > 1000000 {
			fmt.Println("No valid nonce found in 1 million tries")
			break
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Time elapsed: %s\n", elapsed)

}