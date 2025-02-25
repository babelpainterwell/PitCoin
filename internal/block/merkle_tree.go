package block

import "fmt"

// MerklePathItem represents one element in a Merkle proof,
// including the sibling hash and whether that sibling is on the left.
type MerklePathItem struct {
	Hash   [32]byte
	IsLeft bool
}

func (m MerklePathItem) String() string {
    return fmt.Sprintf("MerklePathItem(Hash=%x, IsLeft=%v)", m.Hash, m.IsLeft)
}
