package transaction

import (
	"bytes"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
)

// In terms of tracking orevious outputs, Bitcoin Core keeps a database that stores every UTXO and essential metadata about it.

type TxInput struct {
	OutpointTxID [32]byte // transaction ID where the funding should be spent
	OutputIndex uint32 // index of the output in the transaction
	InputScript []byte // ???
	Sequence uint32 // 
}

type TxOutput struct {
	Amount uint64 // amount of shitcoins
	OutputScript []byte // ???
}

type Transaction struct {
	// with a witness structure for Segregated Witness
	// the marker must be zero 0x00 and the flag must be 0x01 under current P2P protocol
	// marker and flag must not be present for legacy serialization 

	Version uint32
	Marker  uint8
	Flag    uint8
	TxInCount uint32
	TxIn []*TxInput
	TxOutCount uint32
	TxOut []*TxOutput
}

func (tx *Transaction) SerializeTransaction() []byte {
	// serialize the transaction
	var buf bytes.Buffer

	// serialization of all fields
	//
	//
	//
	hashutil.EncodeUint32LE(&buf, tx.Version)
	

	return buf.Bytes()
}


func (tx *Transaction) GetTxID() [32]byte {
	transactionBytes := tx.SerializeTransaction()
	return hashutil.DoubleSha256(transactionBytes)
}
