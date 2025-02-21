package block

import (
	"bytes"

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

func (bh *BlockHeader) SerializeHeader() []byte {
	// serialize the block header

	var buf bytes.Buffer

	// serialization of all fields
	//
	//
	//
	//

	return buf.Bytes()
}

func (bh *BlockHeader) BlockHash() [32]byte {
	headerBytes := bh.SerializeHeader()
	return hashutil.DoubleSha256(headerBytes)
}

type Block struct {
	Header 	 BlockHeader
	Transactions []*transaction.Transaction
}


// ComputeMerkleRoot calculates the merkle root by hashing all TxIDs in the block pairwise 
// If there is an odd number of transactions, re repeat the last one.
// As DoubleSha256 requires a byte slice, we need to convert the [32]byte to a []byte before hashing. 
func (b *Block) ComputeMerkleRoot() [32]byte {

	// 1. Get all the transaction IDs
	txIDs := make([][32]byte, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.GetTxID())
	} 

	if len(txIDs) == 0 {
		return [32]byte{}
	}

	// 2. Convert the transaction IDs to a slice of byte slices
	level := make([][]byte, 0, len(txIDs))
	for _, txID := range txIDs {
		sliceOfTxID := txID[:]
		level = append(level, sliceOfTxID)
	}

	// 3. Build the tree until we are only left with one root 
	for len(level) > 1 {
		var newLevel [][]byte 
		for i := 0; i < len(level); i += 2 {
			// if we have an odd number transactions in current level, we repeat the last one 
			if i + 1 == len(level) {
				newLevel = append(newLevel, doubleSha256Concat(level[i], level[i]))
			} else {
				newLevel = append(newLevel, doubleSha256Concat(level[i], level[i + 1]))
			}
		}
		level = newLevel
	}

	var root [32]byte 
	copy(root[:], level[0])

	return root
}




// to compute a Merkle node, two 32-byte hashes are concatenated and hashed together.
func doubleSha256Concat(first, second []byte) []byte {
	concat := append(first, second...)
	result := hashutil.DoubleSha256(concat)
	return result[:]
}