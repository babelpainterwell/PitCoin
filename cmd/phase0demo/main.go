package main

import (
	"fmt"
	"time"

	"github.com/babelpainterwell/shitcoin/internal/block"
)

func main() {
   
	var blockHeader block.BlockHeader

	blockHeader.Version = 1
	blockHeader.Timestamp = uint32(time.Now().Unix())

	// serialize it using little-endian encoding
	serializedHeader := blockHeader.SerializeHeader()
	fmt.Println(serializedHeader)

	// serialize it using big-endian encoding
	// serializedHeaderBE := blockHeader.SerializeHeaderBE()
	// fmt.Println(serializedHeaderBE)

}
