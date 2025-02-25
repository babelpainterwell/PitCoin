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

	// 2. 
	level := make([][32]byte, 0, len(txIDs))
	level = append(level, txIDs...)

	// 3. Build the tree until we are only left with one root 
	for len(level) > 1 {
		var newLevel [][32]byte 
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

	return level[0]
}

func (b *Block) UpdateMerkleRoot() {
	b.Header.MerkleRoot = b.ComputeMerkleRoot()
}

func (b *Block) RequestMerklePath(txIndex uint32) ([]MerklePathItem, error) {
	// the merkle path is used to verify the transaction in the block

	var mp []MerklePathItem

	txIDs := make([][32]byte, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.GetTxID())
	}

	// txIndex out of range
	if txIndex >= uint32(len(txIDs)) {
		return mp, fmt.Errorf("transaction index out of range")
	}

	// if txIndex is zero and only one transaction, the merkle path is empty
	if txIndex == 0 && len(txIDs) == 1 {
		return mp, nil
	}

	level := make([][32]byte, 0, len(txIDs))
	level = append(level, txIDs...)

	// store the sibling nodes of the transaction ID in the merkle path
	for (len(level) > 1) {
		var newLevel [][32]byte
		for i := 0; i < len(level); i += 2 {
			var left, right [32]byte 
			left = level[i]
			if i + 1 < len(level) {
				right = level[i + 1]
			} else {
				// if odd number, duplicate the last one 
				right = left 
			}

			newLevel = append(newLevel, hashutil.DoubleSha256Concat(left, right))

			// if the txIndex is in this pair, record its sibling 
			if txIndex >= uint32(i) && txIndex < uint32(i + 2) {
				if txIndex == uint32(i) {
					// sibling is the right 
					mp = append(mp, MerklePathItem{Hash: right, IsLeft: false})
				} else {
					// sibling is the left
					mp = append(mp, MerklePathItem{Hash: left, IsLeft: true})
				}
			}
		}
		level = newLevel
		txIndex /= 2 
	}

	return mp, nil 
}


func (b *Block) GetTxIndex(txID [32]byte) (uint32, error) {
	// check if the transaction ID is in the block 
	// find its index in the block 
	txIndex := uint32(len(b.Transactions))
	for i, tx := range b.Transactions {
		if tx.GetTxID() == txID {
			txIndex = uint32(i)
		}
	}
	if txIndex == uint32(len(b.Transactions)) {
		return 0, fmt.Errorf("error to find transaction ID in the block")
	}
	return txIndex, nil
	
}

func (b *Block) VerifyTransaction(txID [32]byte) (bool, error) {
	// verify the transaction in the block
	txIndex, err := b.GetTxIndex(txID)
	if err != nil {
		return false, fmt.Errorf("error to verify transaction; %v", err)
	}
	
	// request for merkle path from the miner
	merklePath, err := b.RequestMerklePath(txIndex)
	if err != nil {
		return false, err
	}

	// verify the merkle path, they should be equal to the merkle root of the block
	currentHash := txID

	// if the merkle path is empty, the transaction is the only one in the block
	if len(merklePath) == 0 {
		return currentHash == b.Header.MerkleRoot, nil
	}

	for _, item := range merklePath {
		if item.IsLeft {
			currentHash = hashutil.DoubleSha256Concat(item.Hash, currentHash)
		} else {
			currentHash = hashutil.DoubleSha256Concat(currentHash, item.Hash)
		}
	}
	return currentHash == b.Header.MerkleRoot, nil
	
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




