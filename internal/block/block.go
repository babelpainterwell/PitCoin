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

func (b *Block) RequestMerklePath(txIndex uint32) ([][32]byte, error) {
	// request the merkle path from the miner
	// the miner will return the merkle path in the form of a slice of 32-byte arrays
	// the merkle path is the sibling nodes of the transaction ID in the merkle tree
	// the merkle path is used to verify the transaction in the block

	mp := make([][32]byte, 0)
	


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
		fmt.Println("Transaction not found in the block")
		return 0, fmt.Errorf("transaction not found in the block")
	}
	return txIndex, nil
	
}

func (b *Block) VeirifyTransaction(txID [32]byte) (bool, error) {
	// verify the transaction in the block
	txIndex, err := b.GetTxIndex(txID)
	if err != nil {
		return false, err
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

	for i, siblingHash := range merklePath {
		// update the isRightChild flag
		isRightChild := (txIndex & (1 << i)) != 0

		if isRightChild {
			currentHash = hashutil.DoubleSha256Concat(siblingHash, currentHash)
		} else {
			currentHash = hashutil.DoubleSha256Concat(currentHash, siblingHash)
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




