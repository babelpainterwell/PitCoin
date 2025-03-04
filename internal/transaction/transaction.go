package transaction

import (
	"bytes"
	"encoding/binary"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
)

// In terms of tracking orevious outputs, Bitcoin Core keeps a database that stores every UTXO and essential metadata about it.


type TxInput struct {
	OutpointTxID 	[32]byte // transaction ID where the funding should be spent
	OutputIndex 	uint32 // index of the output in the transaction
	InputScript 	[]byte // ??? scriptSig (for legacy or P2SH spends; empty for native segwit)
	Sequence 		uint32 // sequence number, as a relative timelock or still for replacement signalling?
	Witness   		[][]byte // witeness stack
}


// don't we need to track the output index?

type TxOutput struct {
	Amount 			uint64 // amount of shitcoins
	OutputScript 	[]byte // ???scriptPubkey
}

type Transaction struct {
	// with a witness structure for Segregated Witness
	// the marker must be zero 0x00 and the flag must be 0x01 under current P2P protocol
	// marker and flag must not be present for legacy serialization 

	Version uint32
	Marker  uint8
	Flag    uint8
	
	TxIns 	[]*TxInput
	TxOuts 	[]*TxOutput

	LockTime uint32 // lock time field
}


func writeVarInt(buf *bytes.Buffer, value uint64) {
	switch {
	case value < 0xfd:
		buf.WriteByte(byte(value))
	case value <= 0xffff:
		buf.WriteByte(0xfd)
		hashutil.EncodeUint16LE(buf, uint16(value))
	case value <= 0xffffffff:
		buf.WriteByte(0xfe)
		hashutil.EncodeUint32LE(buf, uint32(value))
	default:
		buf.WriteByte(0xff)
		hashutil.EncodeUint64LE(buf, value)
	}
}

// write a variable length byte slice to the buffer prefixed with a compactSize length
func writeVarBytes(buf *bytes.Buffer, data []byte) {
	writeVarInt(buf, uint64(len(data)))
	buf.Write(data)
}

// returns true if any input contains witness data
func (tx *Transaction) isSegwit() bool {
	for _, in := range tx.TxIns {
		if len(in.Witness) > 0 {
			return true
		}
	}
	return false
}



func (tx *Transaction) SerializeTransaction() []byte {
	// serialize the transaction
	var buf bytes.Buffer

	// Write the version 
	binary.Write(&buf, binary.LittleEndian, tx.Version)

	isSegwit := tx.isSegwit()

	if isSegwit {
		// write marker and flag for segwit
		// per current P2P protocol, marker must be 0x00 and flag must be 0x01
		if tx.Marker != 0x00 || tx.Flag != 0x01 {
			panic("invalid marker and flag for segwit transaction")
		}
		buf.WriteByte(tx.Marker)
		buf.WriteByte(tx.Flag)
	}
	

	return buf.Bytes()
}


func (tx *Transaction) GetTxID() [32]byte {
	transactionBytes := tx.SerializeTransaction()
	return hashutil.DoubleSha256(transactionBytes)
}
