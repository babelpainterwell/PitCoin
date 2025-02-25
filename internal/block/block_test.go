package block

import (
	"testing"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
	"github.com/babelpainterwell/shitcoin/internal/transaction"
)

func createTestExamples() (*Block, [][32]byte) {
	tx1 := &transaction.Transaction{Version: 1}
	tx2 := &transaction.Transaction{Version: 2}
	tx3 := &transaction.Transaction{Version: 3}
	tx4 := &transaction.Transaction{Version: 4}

	tx1ID := tx1.GetTxID()
	tx2ID := tx2.GetTxID()
	tx3ID := tx3.GetTxID()
	tx4ID := tx4.GetTxID()

	blk := &Block{
		Header: BlockHeader{
			Version: 1,
		},
		Transactions: []*transaction.Transaction{tx1, tx2, tx3, tx4},
	}

	return blk, [][32]byte{tx1ID, tx2ID, tx3ID, tx4ID}
}

func TestUpdateMerkleRoot(t *testing.T) {
	
	blk, txIDs := createTestExamples()

	left := hashutil.DoubleSha256Concat(txIDs[0], txIDs[1])
	right := hashutil.DoubleSha256Concat(txIDs[2], txIDs[3])
	expected_root := hashutil.DoubleSha256Concat(left, right)

	blk.UpdateMerkleRoot()

	if blk.Header.MerkleRoot != expected_root {
		t.Errorf("Expected Merkle root to be %x, but got %x", expected_root, blk.Header.MerkleRoot)
	}

}

// TestVerifyTransaction verifies that the VerifyTransaction function correctly rebuilds
// the Merkle root using the Merkle path.
func TestVerifyTransaction(t *testing.T) {
	
	blk, txIDs := createTestExamples()

	blk.UpdateMerkleRoot()

	// Verify transaction 3
	valid, err := blk.VerifyTransaction(txIDs[2])
	if err != nil {
		t.Fatalf("VerifyTransaction() failed: %v", err)
	}
	if !valid {
		t.Errorf("Expected VerifyTransaction to succeed for tx3, but it failed")
	}
}

func TestRequestMerklePath(t *testing.T) {

	blk, txIDs := createTestExamples()

	blk.UpdateMerkleRoot()

	left := hashutil.DoubleSha256Concat(txIDs[0], txIDs[1])

	// Request Merkle path for transaction 3
	mp, err := blk.RequestMerklePath(2)
	if err != nil {
		t.Fatalf("RequestMerklePath() failed: %v", err)
	}

	if len(mp) != 2 {
		t.Errorf("Expected Merkle path length to be 2, but got %d", len(mp))
	}

	if mp[0].Hash != txIDs[3] {
		t.Errorf("Expected right sibling hash to be %x, but got %x", txIDs[3], mp[0].Hash)
	}

	if mp[1].Hash != left {
		t.Errorf("Expected left sibling hash to be %x, but got %x", left, mp[1].Hash)
	}

}