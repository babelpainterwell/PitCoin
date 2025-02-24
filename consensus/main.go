package main

// a minimal demo of proof of work consensus

import (
	"time"

	"github.com/babelpainterwell/shitcoin/internal/block"
	"github.com/babelpainterwell/shitcoin/internal/transaction"
)

func main() {
	b := &block.Block{
		Header: block.BlockHeader{
			Version: 1,
			Timestamp: uint32(time.Now().Unix()),
			Target: [32]byte{0x00, 0x00, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			Nonce: 0,
		},
		Transactions: []*transaction.Transaction{},
	}

	// update the merkle root of the block
	b.ComputeMerkleRoot()

	// mine the block
	b.Mine()
}