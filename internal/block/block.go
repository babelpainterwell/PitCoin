package block

type BlockHeaser struct {
	Version       int32
	PrevBlockHash [32]byte
	MerkleRoot    [32]byte
	Timestamp     uint32
	Bits          uint32
	Nonce         uint32
}