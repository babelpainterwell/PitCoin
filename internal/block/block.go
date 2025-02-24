package block

import (
	"bytes"
	"fmt"
	"time"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
	"github.com/babelpainterwell/shitcoin/internal/transaction"
)

// Each block contains a cyrptographic hash of the previous block, a timestamp, and transaction data.


type BlockHeader struct {
	// 80 bytes in total
	// to calculate the target, with the first two hex digits for the exponent and the rest for the coefficient 
	// target = coefficient * 2^(8*(exponent-3))

	Version       uint32 
	PrevBlockHash [32]byte 
	MerkleRoot    [32]byte 
	Timestamp     uint32 
	Target        [32]byte // 256 bits or uint32 as Andreas Antonopoulos suggests??? 
	Nonce         uint32 
}

func (bh *BlockHeader) SerializeHeader() []byte {

	var buf bytes.Buffer

	// serialization of all fields
	hashutil.EncodeUint32LE(&buf, bh.Version)
	buf.Write(bh.PrevBlockHash[:])
	buf.Write(bh.MerkleRoot[:])
	hashutil.EncodeUint32LE(&buf, bh.Timestamp)
	buf.Write(bh.Target[:])
	hashutil.EncodeUint32LE(&buf, bh.Nonce)

	return buf.Bytes()
}


func (bh *BlockHeader) HashBlock() [32]byte {
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
				newLevel = append(newLevel, hashutil.DoubleSha256Concat(level[i], level[i]))
			} else {
				newLevel = append(newLevel, hashutil.DoubleSha256Concat(level[i], level[i + 1]))
			}
		}
		level = newLevel
	}

	var root [32]byte 
	copy(root[:], level[0])

	return root
}

func (b *Block) UpdateMerkleRoot() {
	b.Header.MerkleRoot = b.ComputeMerkleRoot()
}

func (b *Block) Mine() {
	// the mining is successful once a nonce making the hash of the block header less than the target is found 
	// Given the block processing speed in the real world, we set a time limit of 600 seconds for mining

	startTime := time.Now()

	for time.Since(startTime) < 600 * time.Second {

		currHeaderHash := b.Header.HashBlock()

		if bytes.Compare(currHeaderHash[:], b.Header.Target[:]) < 0 {
			fmt.Println("Found a valid nonce")
			fmt.Println("Nonce: ", b.Header.Nonce)
			fmt.Println("Hash (little endian): ", currHeaderHash)
			fmt.Println("Took time: ", time.Since(startTime))
			return
		}

		b.Header.Nonce++

		// to prevent nonce overflow
		if b.Header.Nonce == 0 {
			fmt.Println("Nonce overflow. Mining failed.")
			return
		}
	}

	fmt.Println("Mining failed Due to Time Limit -- 600 seconds")

}



