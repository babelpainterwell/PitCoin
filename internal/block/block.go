package block

import (
	"github.com/babelpainterwell/shitcoin/internal/hashutil"
	"github.com/babelpainterwell/shitcoin/internal/transaction"
)

// Each block contains a cyrptographic hash of the previous block, a timestamp, and transaction data.


type BlockHeader struct {
	// 80 bytes in total

	Version       uint32 
	PrevBlockHash [32]byte 
	MerkleRoot    [32]byte 
	Timestamp     uint32 
	Target        uint32 
	Nonce         uint32 
}

type Block struct {
	Header 	 BlockHeader
	Transactions []*transaction.Transaction
}

func (b *Block) ComputeMerkleRoot() [32]byte {
	//
}




// to compute a Merkle node, two 32-byte hashes are concatenated and hashed together.
func doubleSha256Concat(first, second [32]byte) [32]byte {
	concat := append(first[:], second[:]...)
	return hashutil.DoubleSha256(concat)
}